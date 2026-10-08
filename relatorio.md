# Relatório de testes — Keep v0.2.0

- **Commit testado:** `154c5e3` (branch `main`, árvore limpa)
- **Ambiente:** macOS (APFS, *case-insensitive*), Go 1.26.0
- **Data:** 2026-10-08

## Resumo

| Severidade | Qtd. | Significado |
|---|---|---|
| Crítico | 3 | perda ou corrupção silenciosa de dados em uso realista |
| Alto | 7 | dado errado devolvido, garantia documentada quebrada ou risco de segurança |
| Médio | 9 | bug real, mas que exige um gatilho menos comum, ou uma lacuna séria de API |
| Baixo | 13 | caso de borda, inconsistência, ergonomia |
| Docs / Info | 6 | documentação errada, desempenho, observações |

### Como foi testado

Um módulo Go isolado (`replace` apontando para este repositório, sem tocar em nenhum
arquivo do projeto), com cerca de 45 testes cobrindo os dois backends (`memory` e `file`):

- **Injeção de falhas:** um wrapper de `StorageDeps` que faz um `Write`, `Delete` ou `Read`
  específico falhar, para exercitar os caminhos de erro `Internal` e as ordens de escrita
  da DenseRecordPattern.
- **Fuzz baseado em modelo:** 10 sementes × 3000 operações aleatórias em memória e 1500 em
  arquivo (insert com duplicatas que só mudam de caixa, remove, update de `Key` incluindo
  troca de caixa, update simples, insert/remove aninhado). Os invariantes foram conferidos
  contra um modelo de referência a cada 50 operações.
- **Concorrência:** escritores concorrentes, e um único escritor com leitores concorrentes,
  além de `go test -race`.
- **Adaptadores direto no contrato:** `filestorage` e `memstorage` chamados sem passar pela
  biblioteca.
- **Validação de schema, tipos, Unicode, limites, paginação, evolução de schema e
  desempenho.**

### O que funciona bem

- Os 14 exemplos passam (`agnos run-examples`); `go build` e `go vet` passam sem avisos.
- **Com um único escritor e sem falhas de storage, o núcleo está correto.** O fuzz rodou
  31.500 operações sem nenhuma divergência do modelo: listagem, valores, `FindByKey`,
  aninhados, ids monotônicos e sem reuso, unicidade sem diferenciar maiúsculas.
- Remover todos os registros não deixa vazamento: sobram só `size` e `last-id`, como
  esperado, nos dois backends.
- `go test -race` não acusa nenhuma *data race* do Go. Os problemas de concorrência abaixo
  são lógicos, não de memória.

A maior parte dos bugs graves está em três pontos: **(a)** o que acontece quando uma
operação falha no meio, **(b)** concorrência e **(c)** o adaptador `filestorage`.

---

## Crítico

### 1. Escritores concorrentes corrompem o banco sem reportar erro

**Onde:** `sandbox/internal/record/record.go` (`Insert`, `Update`, `Remove`) — nenhuma
operação usa trava, `WriteIfAbsent` nem `WriteIfValueEquals`.

**Reprodução:** 16 goroutines × 25 `Insert` com emails distintos, mais 16 goroutines
inserindo o mesmo email ao mesmo tempo.

**Observado:**
- `memory`: as **400 chamadas retornam sucesso**, mas `ListAll` lista só 50–77 registros, com
  2–6 **ids duplicados** na lista. Só 172–228 dos 400 `FindByKey` devolvem o registro certo.
- `file`: 145 de 400 retornam sucesso e `ListAll` lista 37. **9 de 16 inserts simultâneos
  do mesmo valor `Key` (que deveria ser único) retornam sucesso.**

**Por que é crítico:** em Go, um handler HTTP já é concorrente por padrão. A única menção a
"single writer" está em `docs/DenseRecordPattern`. README, Schemas, Errors e os
doc-comments da API (`Collection.Insert`) não avisam, e nada impede o uso concorrente.

**Correção sugerida:** no mínimo, um `sync.Mutex` por prefixo de coleção dentro do processo.
Para múltiplos processos: `WriteIfAbsent` para reivindicar a entrada de índice e o slot da
lista, e `WriteIfValueEquals` como CAS em `last-id` e `size` (o contrato já oferece os dois),
ou `Lock` no prefixo. E documentar no README e nos doc-comments.

### 2. Um `Insert` que falha com `Internal` deixa um registro "meio visível" que depois apaga outro registro

**Onde:** `record.go:423` grava `position` e `record.go:428-440` grava o índice **antes** do
ponto de commit (`size`, `record.go:446`). `ResolveByID` (`record.go:458`) considera vivo
qualquer id que tenha a chave `position`.

**Reprodução:** um storage cujo `Write` em `db/user/size` falha uma vez.

```go
users.Insert({"email": "ghost@x.com"})   // -> Internal
users.FindByKey("email", "ghost@x.com")  // -> ok == true   (!)
users.FindByID(1)                         // -> ok == true   (!)
users.Insert({"email": "ghost@x.com"})   // -> KeyConflict, para sempre
legit := users.Insert({"email": "legit@x.com"})  // ocupa a posição 1
ghost.Remove()                            // -> legit some do ListAll, mas FindByID(legit) continua ok
```

**Contradiz:** `docs/DenseRecordPattern:65`: *"never a half-visible record"*.

**Correção sugerida:** definir "vivo" como `position(id)` existe **e** `position <= size`
**e** `list[position] == id`, e aplicar isso em `ResolveByID`, `Remove`, `Update` e
`InsertNested`. No `Insert`, tratar como livre uma entrada de índice cujo dono não está vivo.

### 3. Repetir um `Remove` que falhou tira outro registro da listagem, em cascata

**Onde:** `record.go:157-197`. O `Remove` confia em `position(id)` sem verificar
`list[position] == id`, e a chave `position` só é apagada no fim (`record.go:232`).

**Reprodução:** com `[a, b]`, o `Delete` da entrada de índice falha, então `a.Remove()`
retorna `Internal`. Nesse ponto `b` já foi movido para a posição 1 e `size = 1`, mas
`a.position` continua 1.

```go
a.Remove()          // retry: retorna nil
users.ListAll()     // -> []   e b continua vivo (FindByID ok)
c := users.Insert(...)
b.Remove()          // -> c também some do ListAll, embora continue vivo
```

Repetir após `Internal` é a reação natural de quem chama, e a doc garante que remover um
registro que já sumiu é seguro.

**Correção sugerida:** se `p > size` ou `list[p] != id`, o registro já foi desligado da
lista. Basta apagar as chaves dele, sem mexer em `list` nem em `size`.

---

## Alto

### 4. Um handle de `Record` antigo continua escrevendo depois do `Remove`: zumbi, chave bloqueada para sempre e vazamento

**Onde:** `UpdateFactory` (`record.go:79`) e `InsertNestedFactory` (`record.go:271`) não
verificam se o registro está vivo. A checagem de unicidade do `Insert` (`record.go:390`) usa
`Exists` no índice sem olhar se o dono está vivo.

**Reprodução:**

```go
a.Remove()
a.Update("email", "z@x.com")  // -> nil
a.Update("age", 33)           // -> nil
a.Get("age")                  // -> 33  (registro zumbi)
users.Insert({"email": "z@x.com"})  // -> KeyConflict, sem nenhum registro vivo com esse valor
a.InsertNested("sessions", ...)     // -> sucesso
a.Remove()                          // -> no-op: 9 chaves ficam órfãs para sempre
```

Isso é comum na prática: guardar o resultado de `ListAll`, remover alguns e depois atualizar.

**Contradiz:** `docs/Errors`: `KeyConflict` significa *"another **live** record already holds
that value"*.

**Correção sugerida:** checar se o registro está vivo no início de `Update` e `InsertNested`
e retornar erro. Tratar como livre a entrada de índice cujo dono está morto.

### 5. Um `Update` de `Key` que falha deixa um índice "mentiroso": `FindByKey` devolve um registro com outro valor

**Onde:** `FindByKeyFactory` (`sandbox/internal/collection/collection.go:27-42`) segue a
entrada de índice sem conferir o invariante 3 ("toda entrada aponta para um registro cujo
valor atual tem esse hash").

**Reprodução:** a escrita de `values/email` falha durante `a.Update("email", "new@x.com")`.

```go
users.FindByKey("email", "new@x.com")  // -> devolve `a`, e a.Get("email") == "old@x.com"
users.Insert({"email": "new@x.com"})   // -> KeyConflict
```

**Contradiz:** o comentário de `record.go:73-78` diz que a entrada velha *"resolves to nothing
and is overwritten by the next write"*. Na prática ela resolve para o registro errado e bloqueia
a próxima escrita.

**Correção sugerida:** `FindByKey` deve confirmar que `HashIndexValue(Get(field))` bate com
o hash da entrada. `Insert` e `Update` devem tratar uma entrada divergente como livre.

### 6. "Readers are always safe" é falso: com um único escritor, leitores recebem registros fantasma (id 0) e erros

**Onde:** `ListRange` (`record.go:479`) lê `size` uma vez e depois os slots. Um slot apagado
nesse meio-tempo cai em `ReadInt` (`dense.go:279-285`), que trata ausência como `0`, e vira
`Build(id 0)`. No `filestorage`, `os.WriteFile` (`filestorage.go:86`) trunca o arquivo antes
de escrever, então um leitor pode ver `""`.

**Observado:**
- `memory`: 3000 registros, uma goroutine removendo e 4 fazendo `ListAll` → **5380 registros
  com `ID == 0`**.
- `file`: um escritor (insert, remove e update) por 1,5 s e 8 leitores → **1185 erros
  `Internal` e 278 registros fantasma**.

**Contradiz:** `docs/DenseRecordPattern:123`: *"Keep assumes a single writer. Readers are
always safe"*.

**Correção sugerida:** ver itens 7 e 13. Um slot ausente nunca deve virar id 0, e a escrita
em arquivo precisa ser atômica.

### 7. `filestorage`: a escrita não é atômica, e um crash no meio deixa um arquivo vazio que inutiliza a coleção

**Onde:** `filestorage.go:86` (`os.WriteFile` trunca e depois escreve) e
`filestorage.go:119,158`. Também não há `fsync`.

**Reprodução:** truncar `db/user/size`, que é o estado que um crash entre o truncamento e a
escrita deixa. A partir daí, **todo** `Insert` e `ListAll` falha com
`strconv.ParseInt: parsing "": invalid syntax`, e não existe ferramenta para recuperar
(item 16).

**Contradiz:** `docs/StorageContract:83`: *"a crash mid-sequence leaves debris that recovery
deletes, never corrupt data"*. Toda a ordem de escrita da DenseRecordPattern pressupõe que a
escrita de uma chave seja atômica.

**Correção sugerida:** escrever num arquivo temporário no mesmo diretório, fazer `fsync` e
depois `os.Rename`. Opcionalmente, `fsync` também no diretório.

### 8. `Link` aceita um `Record` de outra coleção e resolve para o registro errado

**Onde:** `dense.go:211`. Do `api.Record` só se lê o `ID`, nunca o `Prefix`.

**Reprodução:**

```go
alice := users.Insert({"email": "alice@x.com"})  // user id 1
p1    := posts.Insert({"slug": "p1"})            // post id 1
bob   := users.Insert({"email": "bob", "best": p1})  // aceito, e p1 é um POST
bob.GetLink("best")                              // -> alice
posts.Insert({"slug": "p2", "author": session})  // registro aninhado também é aceito
```

**Contradiz:** `docs/Schemas:138`: *"What it cannot do is resolve to the wrong record"*.

**Correção sugerida:** quando o valor for um `api.Record`, comparar o `Prefix` dele com o
prefixo do `Target` (o resolver já está disponível no caminho de insert e update) e retornar
`InvalidField` se divergir.

### 9. `filestorage` num FS sem diferenciação de maiúsculas (padrão do macOS e do Windows): campos e schemas que diferem só na caixa se sobrescrevem

**Onde:** `filestorage.go:30-40`. `url.PathEscape` preserva a caixa, e o FS a ignora.

**Reprodução (nesta máquina):**
- Os campos `"Name"` e `"name"` do mesmo schema, inseridos com `"DISPLAY"` e `"login"`,
  passam a ler os dois `"login"`.
- Os schemas `"User"` e `"user"` compartilham os mesmos registros.
- No backend `memory` o comportamento é correto, então o mesmo programa se comporta diferente
  dependendo do backend.

**Contradiz:** o contrato exige *"the flattening is injective"*.

**Correção sugerida:** uma codificação injetiva que não dependa da caixa (por exemplo, escapar
letras maiúsculas como `%XX`, ou hex/base32 do segmento), ou rejeitar esses nomes na
validação.

### 10. `filestorage`: os segmentos `.`, `..` e `""` não são escapados, o que permite escrever fora do diretório base e gera colisões

**Onde:** `filestorage.go:30-40`. `url.PathEscape` não escapa `.`, segmentos vazios são
descartados e `filepath.Join` normaliza `..`.

**Reprodução:**
- `filestorage.New(base).Write(["..", "escaped.txt"])` grava **no diretório pai** de `base`.
- `Props.Path: "../../outside"` grava fora da base.
- Um schema chamado `".."` sob `Path: "db"` grava `base/size`.
- `["a", ".", "b"]` e `["a", "b"]` são o mesmo arquivo, e `["x", "", "y"]` é igual a `["x", "y"]`.
- Os campos `"."` e `""` viram o mesmo arquivo. Depois de gravar o campo `"."`, nenhum outro
  campo do registro pode ser gravado: `mkdir …/values: not a directory`.

**Contradiz:** o doc-comment de `filestorage.go:25-29` (*"without ever escaping the base
directory"*) e `docs/StorageContract:31`. Vira um problema de segurança se nomes de
schema ou o `Path` vierem de configuração ou de entrada externa.

**Correção sugerida:** escapar explicitamente `.`, `..` e `""` (ou prefixar todo segmento
com um caractere fixo) e validar o `Path`.

---

## Médio

### 11. Os nomes internos do layout não são reservados: um campo aninhado chamado `values` ou `position` colide com os dados do registro

**Onde:** `dense.go:75-95`. `SubPrefix = {c}/{id}/{field}` divide o namespace com
`{c}/{id}/position` e `{c}/{id}/values/…`.

**Reprodução (nos dois backends):** um schema com `{size: Int, values: Nested}`, inserindo
`size = 3`.
- `ListNested("values")` devolve **3 registros fantasma `[0 0 0]`**, porque lê `values/size`.
- `InsertNested("values", …)` **sobrescreve o campo `size` para 4**.
- Um campo aninhado chamado `"position"` falha no `file` (`not a directory`) e funciona no
  `memory`.

**Correção sugerida:** pôr os aninhados sob um segmento dedicado (`{c}/{id}/nested/{field}`,
o que muda o layout e exige migração) ou rejeitar nomes reservados na validação de schema
(item 14).

### 12. Erros de storage viram "não encontrado"

**Onde:** `FindByKey`, `FindByID`, `GetLink`, `HasValues` e `ListNested` retornam só
`bool` ou `nil`.

**Reprodução:** com o backend falhando em toda leitura, todas as cinco respondem
"ausente". Uma aplicação trata uma queda do storage como inexistência (por exemplo, "usuário
não encontrado, então cria"). `ListNested` não distingue entre "campo não existe", "vazio" e
"erro".

**Correção sugerida:** assinaturas `(Record, bool, *Error)` ou equivalente, e
`ListNested` retornando `([]Record, *Error)`.

### 13. `ListAll`/`List` devolvem um registro fantasma com id 0 quando falta um slot, em vez de erro

**Onde:** `ListRange` usa `ReadInt` (`dense.go:285` trata ausência como 0) para os slots
da lista.

**Reprodução:** apagar `db/user/list/2` e chamar `ListAll`, que devolve `[1 0]` sem erro. No
mesmo estado, o `Remove` retorna `Internal "position 2 of the list is missing"`, ou seja, o
comportamento é inconsistente. É também a causa dos fantasmas do item 6.

**Correção sugerida:** ler os slots com `Read` e checar `found`, retornando `Internal`
(ou pulando o slot, se for para tolerar leitura concorrente).

### 14. Não há validação de schema

**Onde:** `Databases.New` não tem retorno de erro.

**Aceito sem reclamar:**
- nomes de campo duplicados;
- schemas com o mesmo nome (o segundo é ignorado em silêncio);
- `Type: 99`;
- `Link` com `Target` inexistente;
- nomes reservados (item 11);
- nomes `"."`, `""` e `".."` (item 10).

**Reprodução:** um campo `Required` com `Type: 99` torna a coleção **impossível de usar**.
Se o campo é passado, o insert falha com `InvalidField`; se é omitido, falha com
`MissingField`. Um `Link` com `Target: "nope"` (um erro de digitação) é aceito no insert, e
`GetLink` sempre devolve `false`.

**Correção sugerida:** validar em `New` (retornando `(Database, *Error)`) ou expor
`Databases.Validate(props) *Error`.

### 15. O valor zero de `Field.Type` é `Key`, então esquecer o `Type` cria um índice único

**Onde:** `sandbox/api/databases.go:20` (`Key = iota`) e `:57` (`KeyConflict = iota`).

**Reprodução:** `{Name: "city"}` sem `Type` vira uma chave única sem diferenciar maiúsculas,
então inserir "Lisboa" e depois "lisboa" gera `KeyConflict`. Além disso, um `api.Error{}`
zerado se lê como `KeyConflict`.

**Correção sugerida:** começar os `iota` em 1 e tratar 0 como inválido na validação. Isso
muda os valores das constantes, o que é aceitável antes da 1.0.

### 16. Não há ferramenta de recuperação, e a tabela de "debris" da doc está incompleta

**Onde:** `docs/DenseRecordPattern`, seção *Recovery*.

A doc descreve resíduos *"safe to delete"*, mas a API não tem nada para detectá-los ou
removê-los. Além disso, os estados produzidos pelos itens 2, 3 e 5 não se encaixam na
tabela: uma entrada de índice que aponta para um registro com `position` mas com outro valor,
e um registro com `position` cujo slot pertence a outro. Junto com o item 7, um crash exige
cirurgia manual nos arquivos.

**Correção sugerida:** um `Collection.Check()`/`Repair()` que percorra as posições `1..size`
e os ids `1..last-id`. Isso continua sem listagem, só com leituras pontuais.

### 17. A evolução de schema não é suportada e falha em silêncio

**Reprodução:**
- **De `String` para `Key`:** os registros antigos não são indexados. `FindByKey` não os
  encontra, e um insert duplicado é aceito.
- **De `String` para `Int`:** `Get` passa a retornar `Internal "parsing \"old\""`.
- **Campo removido do schema:** o `Remove` deixa o valor dele no storage para sempre.

Não há reindexação nem migração.

### 18. As coleções aninhadas são de segunda classe

Não existem `FindByKey`, `FindByID` nem `List` paginado para aninhados. O índice único
aninhado é "somente escrita": ele impede duplicatas, mas não pode ser consultado. Achar uma
sessão pelo token exige varrer `ListNested`. Um `Link` não pode apontar para um registro
aninhado.

**Contradiz:** `docs/Schemas`: *"It behaves like a top-level collection in every way"*.

### 19. Um `Stringer` nulo causa `panic` dentro da biblioteca

**Onde:** `dense.go:168-169`.

**Reprodução:** `Insert({"email": (*T)(nil)})`, onde `*T` tem um método `String()` que
desreferencia o receptor. O `panic: nil pointer dereference` sai de dentro do Keep e derruba
quem chamou.

**Correção sugerida:** `defer`/`recover` em volta de `String()` (`recover` é builtin, então
cabe no sandbox), convertendo para `InvalidField`.

---

## Baixo

| # | Problema | Evidência |
|---|---|---|
| 20 | **`List` com overflow e parâmetros silenciosos.** `from+chunk-1` estoura em `record.go:488` | `List(2, math.MaxInt)` devolve 0 registros, quando o esperado é 4. `List(1, -1)` devolve tudo. `List(0, 2)` e `List(-7, 2)` passam a começar em 1 sem avisar |
| 21 | **`Lock`/`Unlock` do `filestorage`:** corrida ao assumir um lease expirado (`filestorage.go:225-237`), o sufixo `.keeplock` pode ser produzido por um segmento normal (`:44-45`), `Unlock` não tem dono e um conteúdo ilegível conta como expirado | 3 de 32 goroutines adquiriram o mesmo lease. `Lock(["x"])` sobrescreveu o valor da chave `["x.keeplock"]` com um timestamp, e `Unlock(["x"])` apagou essa chave. O Keep não chama `Lock`, por isso a severidade é baixa |
| 22 | **`ReadAt` entra em `panic`** nos dois adaptadores (`filestorage.go:201`, `memstorage.go:135`) | `ReadAt(k, 1, -3)` e `ReadAt(k, 1, math.MaxInt64)` dão `makeslice: len out of range` / `slice bounds out of range` |
| 23 | **`filestorage` diverge do contrato e do `memory`** | `Exists(["a"])` é `true` quando só `["a","b"]` existe (é um diretório). `Read(["a"])` retorna o erro `is a directory` em vez de `found=false`. As chaves `["a"]` e `["a","b"]` não podem coexistir |
| 24 | **Limites de nome e caminho só no `filestorage`**, sem documentação | Um nome de campo com 50 `é` (100 bytes, que o escape triplica para 300) dá `file name too long`. O aninhamento falha na profundidade 35 por caminho longo demais. O `memory` aceita os dois |
| 25 | **`Int` rejeita tipos inteiros comuns** | `float64(27)`, que é o que `encoding/json` produz, além de `uint`, `uint8..64`, `int8` e `int16`, são rejeitados. `Float` aceita `int`, mas não `uint` |
| 26 | **Unicode nas `Key`:** `ToLower` não é *case folding* nem normalização (`dense.go:106`) | `"ΟΔΟΣ"`/`"οδος"`, `"café"` NFC/NFD e `"STRASSE"`/`"straße"` são aceitas como chaves diferentes (emails homógrafos). Já `"Keep"` (sinal de Kelvin) colide com `"keep"` |
| 27 | **`Error.Field` não é determinístico** (iteração de map em `record.go:349`) | Um insert com 3 campos inválidos, repetido 200 vezes, relatou `nope1` 143×, `nope2` 31× e `nope3` 26× |
| 28 | **A paginação pula registros quando há um `Remove` entre páginas** | O registro 20 nunca aparece. Isso está documentado na DenseRecordPattern, mas não no doc-comment de `List` nem em Schemas |
| 29 | **`Collection.Fields` compartilha o array com `Props`** | `users.Fields[1].Type = api.Key` muda o schema do banco inteiro e o `Props` de quem chamou |
| 30 | **Não dá para limpar um campo opcional** | `Update("bio", nil)` retorna `InvalidField`. Um campo `Required` com `nil` dá `InvalidField`, não `MissingField` |
| 31 | **`Record.String` ambíguo** | Um registro sem valores imprime `"{id: 1, }"`. `a = "x, b: y"` imprime igual a `a = "x", b = "y"` |
| 32 | **`HasValues` permissivo** | `HasValues(nil)` é `true`, inclusive num registro removido. Nomes desconhecidos ou aninhados retornam `false` em vez de erro |

---

## Documentação e informativos

| # | Item | Detalhe |
|---|---|---|
| 33 | **Desempenho do `filestorage`** | Com 3000 registros: insert 1,6 ms/op, remove 2,0 ms/op, `FindByKey` 49 µs, `ListAll+Get` 70 µs por registro. No `memory`: 10 µs, 8 µs, 2,4 µs e 2,4 µs. Cada `Write` faz `MkdirAll` e cada `Delete` tenta podar até a profundidade inteira |
| 34 | **Permissões `0644`/`0755`** | Os arquivos do banco (emails, `Bytes`) ficam legíveis por todos os usuários da máquina. Sugestão: `0600`/`0700`, ou permissões configuráveis |
| 35 | **O binding `standard` usa `"."`** (`filestorage.go:72`) | O local do banco depende do diretório de trabalho do processo |
| 36 | **O título do README renderiza `# <no value>`** | Variável de template não preenchida na geração |
| 37 | **O README diz coisas desatualizadas** | *"eleven runnable programs"*, mas existem 14. *"the library only ever calls the eleven single-key functions"*, enquanto `docs/StorageContract` diz que o Keep só chama 4 (`Write`, `Read`, `Exists`, `Delete`) |
| 38 | **O exemplo do README ignora o `ok` de `Collection("user")`** | Um erro de digitação no nome devolve um `Collection{}` zerado, e o `Insert` dá `panic` por função nil |

---

## Prioridade sugerida de correção

1. **Definir "vivo" corretamente** (itens 2, 3, 4) e **validar entradas de índice** (itens 4,
   5). Com uma única função `isLive(prefix, id)` (`position` existe, `position <= size` e
   `list[position] == id`) mais a regra "uma entrada de índice cujo dono não está vivo ou
   cujo valor diverge está livre", os itens 2–5 se resolvem e a tabela de *Recovery* da doc
   passa a ser verdade.
2. **Escrita atômica no `filestorage`** e **nunca transformar um slot ausente em id 0**
   (itens 6, 7, 13).
3. **Concorrência:** um mutex por coleção no processo, já, e documentar em destaque (item 1).
4. **Codificação de segmentos no `filestorage`** (itens 9, 10, 21, 23) e **validação de
   schema** com nomes reservados (itens 11, 14, 15).
5. **`Link` checando a coleção** (item 8).
