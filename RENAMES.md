# Proposta de renomeação — nomes do Keep

Escopo: só nomes do Keep. Nada nativo do agnos é tocado (ver [O que fica](#o-que-fica)).

## 1. API pública — `sandbox/api/databases.go` (breaking)

| Atual | Proposto | Por quê |
|---|---|---|
| `Item.Itens`, `Schema.Itens` | `Fields` | "Itens" é português numa API em inglês, e o mesmo dado se chama `Items` em `SchemaItem` e `SchemaInstance` |
| `SchemaItem.Items`, `SchemaInstance.Items` | `Fields` | mesmo dado que `Itens`; um nome só |
| `Item` | `Field` | o próprio comentário diz "Item describes one field"; o resto já fala field (`fieldName`, `InvalidField`, `MissingField`) |
| `SchemaItem` | `Record` | doc e código chamam de "record" (a variável é `record`); "Item" já é usado para campo |
| `SchemaInstance` | `Collection` | a doc chama de "collection"; "Schema" é a declaração, não o objeto vivo |
| const `Database` (tipo de campo) | `Nested` | é uma coleção aninhada, não um banco — e libera o nome `Database` |
| `DatabaseHandle` | `Database` | "Handle" é ruído: `lib.Databases.New(props)` devolve um `Database` |
| `DatabaseHandle.GetSchema` | `Collection` | devolve uma coleção, não um schema: `db.Collection("user")` |
| `SchemaInstance.NewItem` | `Insert` | grava no storage; casa com `Update` / `Remove` |
| `SchemaItem.NewSubItem` | `InsertNested` | idem, para o campo aninhado |
| `SchemaItem.ListAll(field)` | `ListNested` | mesmo nome que `SchemaInstance.ListAll()`, semântica diferente |
| `SchemaItem.CheckKeysPresence(keys)` | `HasValues(fields)` | não tem a ver com campos `Key`: checa se os campos têm valor |
| `Error.Key`, `Error.KeyValue` | `Error.Field`, `Error.Value` | é o nome de um campo de qualquer tipo; "Key" confunde com o tipo `Key` |
| const `NotFound` | `NoValue` | significa "campo sem valor", não "registro não encontrado" (esse caso é `ok == false`) |
| `FindByKey(key, keyValue)` | `FindByKey(field, value)` | só os parâmetros, pelo mesmo motivo de `Error.Key` |

`Item`, `SchemaItem`, `SchemaInstance` e `DatabaseHandle` andam juntos, e `Database` só fica livre
depois que a const `Database` vira `Nested`.

## 2. Internos — `sandbox/internal/` (sem impacto externo)

| Atual | Proposto | Por quê |
|---|---|---|
| pacote `schemainstance` | `collection` | acompanha a API |
| pacote `schemaitem` | `record` | acompanha a API |
| `schemaitem.New` | `Insert` | esse `New` grava no storage; quem só monta é `Build` — hoje o nome engana |
| `databases.NewHandle`, `GetSchemaFactory` | `NewDatabase`, `CollectionFactory` | acompanha a API |
| `dense.FindItem` | `FindField` | acompanha a API |
| `dense.ReadCount` | `ReadInt` | também lê ids (`ListRange`, `ClearCollection`) e é o par de `WriteInt` |
| `schemaitem.ResolveLive(rawId)` | `ResolveRawId` | as duas `Resolve*` só devolvem registros vivos; o que difere é receber bytes crus |

## 3. Contrato de storage e adapters

| Atual | Proposto | Por quê |
|---|---|---|
| `storagedeps.Contract.UnLock` | `Unlock` | grafia do Go (`sync.Mutex.Unlock`) |
| `WriteIfKeyNotExists` | `WriteIfAbsent` | mais curto, e "absent" é o termo que o próprio contrato usa |
| parâmetro `old_value` | `oldValue` | snake_case não é Go |
| binding `native` | `memory` | "native" sugere algo do sistema operacional; é o contrário — memstorage, nada gravado em disco |

## 4. Opcional

| Atual | Proposto | Por quê |
|---|---|---|
| `SchemaItem.Id`, `FindById`, `LastIdKey`, `ParseId` | `ID`, `FindByID`, `LastIDKey`, `ParseID` | convenção de iniciais do Go, cobrada por linters |

## O que fica

- **Nomes do agnos:** `Sandbox`, `Config`, `Deps`, `ProjectSandbox`, `ProjectConfig`, `Constructor`,
  `Bind`, `Contract`, o padrão `New<X>` / `<Campo>Factory`, `stddeps` / `stringsdeps` / `hashdeps`
  com `osstd` / `sha256hash` / `stdstrings`, a binding `standard` e a estrutura de `sandbox/`.
- **Nomes do Keep que estão bons:** `Databases`, `Info`, `Props`, `Schema`, `Key` / `Int` / `Float` /
  `String` / `Link` / `Bytes`, `Target`, `Required`, `dense`, `liberror`, `storagedeps`.
  `filestorage` e `memstorage` já seguem o padrão `<impl><x>`.
- **Segmentos de chave gravados no storage** (`size`, `last-id`, `list`, `position`, `values`,
  `keys`): são o formato em disco, e trocá-los quebraria bancos que já existem.

## Ao aplicar

- **Versão:** as seções 1 e 3 quebram compatibilidade. Com o Keep em `v0.1.0`, cabe uma `v0.2.0`
  (`version:` em `AgnosConfig/project.yaml`, depois `agnos build`).
- **`AgnosConfig/structure.yaml`:** renomear `schemainstance` / `schemaitem` exige atualizar as
  entradas deles, senão `agnos verify` falha.
- **Exemplos e docs:** todos os exemplos, e as docs `Schemas` e `Errors`, usam `Itens` / `NewItem` /
  `GetSchema`. `docs/PublicApi` é regenerado pelo build.
- **Goldens (`result.yaml`):** não devem mudar, porque o formato em disco continua o mesmo.
