package lib_test

import (
	"testing"

	"github.com/MateusMoutinhoOrg/Keep/adapters/native"
	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
	"github.com/MateusMoutinhoOrg/Keep/pkg/lib"
)

var userSchema = lib.Schema{
	Name: "users",
	Itens: []lib.Item{
		{Name: "email", Type: lib.Key, Required: true},
		{Name: "username", Type: lib.Key, Required: true},
		{Name: "age", Type: lib.Int, Required: true},
		{
			Name: "sessions",
			Type: lib.Database,
			Itens: []lib.Item{
				{Name: "token", Type: lib.Key, Required: true},
				{Name: "creation", Type: lib.Int, Required: true},
			},
		},
	},
}

func newUsers(t *testing.T, d deps.Deps) *lib.SchemaInstance {
	t.Helper()
	keep := lib.New(d)
	db := keep.NewDatabase(lib.Props{
		Path:    "testdb/",
		Schemas: []lib.Schema{userSchema},
	})
	users := db.GetSchema("users")
	if users == nil {
		t.Fatal("GetSchema returned nil")
	}
	return users
}

func mustCreate(t *testing.T, users *lib.SchemaInstance, email, username string, age int) *lib.SchemaItem {
	t.Helper()
	item, err := users.NewItem(map[string]any{"email": email, "username": username, "age": age})
	if err != nil {
		t.Fatalf("NewItem(%s): %v", email, err)
	}
	return item
}

func runWithAdapters(t *testing.T, test func(t *testing.T, d deps.Deps)) {
	t.Run("native", func(t *testing.T) {
		test(t, native.New())
	})
	t.Run("standard", func(t *testing.T) {
		test(t, standard.NewWithBase(t.TempDir()))
	})
}

func TestCreateAndFind(t *testing.T) {
	runWithAdapters(t, func(t *testing.T, d deps.Deps) {
		users := newUsers(t, d)
		mustCreate(t, users, "a@x.com", "alice", 30)

		found := users.FindByKey("email", "a@x.com")
		if found == nil {
			t.Fatal("FindByKey by email returned nil")
		}
		age, e := found.Get("age")
		if e != nil || age.(int64) != 30 {
			t.Fatalf("Get(age) = %v, %v; want 30", age, e)
		}
		// Lookup is case-insensitive (normalization before hashing).
		if users.FindByKey("email", "A@X.com") == nil {
			t.Fatal("case-insensitive lookup failed")
		}
		if users.FindByKey("email", "missing@x.com") != nil {
			t.Fatal("lookup of missing record should return nil")
		}
		if !found.CheckKeysPresence([]string{"email", "username", "age"}) {
			t.Fatal("CheckKeysPresence should be true for stored fields")
		}
	})
}

func TestUniquenessConflict(t *testing.T) {
	runWithAdapters(t, func(t *testing.T, d deps.Deps) {
		users := newUsers(t, d)
		mustCreate(t, users, "a@x.com", "alice", 30)
		_, err := users.NewItem(map[string]any{"email": "A@X.com", "username": "other", "age": 1})
		if err == nil || err.Type != lib.KeyConflict {
			t.Fatalf("want KeyConflict, got %v", err)
		}
	})
}

func TestMissingRequiredField(t *testing.T) {
	runWithAdapters(t, func(t *testing.T, d deps.Deps) {
		users := newUsers(t, d)
		_, err := users.NewItem(map[string]any{"email": "a@x.com", "username": "alice"})
		if err == nil || err.Type != lib.MissingField {
			t.Fatalf("want MissingField, got %v", err)
		}
	})
}

func TestUpdateIndexedField(t *testing.T) {
	runWithAdapters(t, func(t *testing.T, d deps.Deps) {
		users := newUsers(t, d)
		mustCreate(t, users, "a@x.com", "alice", 30)
		mustCreate(t, users, "b@x.com", "bob", 40)

		alice := users.FindByKey("email", "a@x.com")
		// Conflict with bob's email must be rejected.
		if e := alice.Update("email", "b@x.com"); e == nil || e.Type != lib.KeyConflict {
			t.Fatalf("want KeyConflict, got %v", e)
		}
		// Legit re-index: old lookup dies, new lookup works.
		if e := alice.Update("email", "new@x.com"); e != nil {
			t.Fatalf("Update(email): %v", e)
		}
		if users.FindByKey("email", "a@x.com") != nil {
			t.Fatal("old email should no longer resolve")
		}
		if users.FindByKey("email", "new@x.com") == nil {
			t.Fatal("new email should resolve")
		}
		// Non-indexed update.
		if e := alice.Update("age", 31); e != nil {
			t.Fatalf("Update(age): %v", e)
		}
		age, _ := alice.Get("age")
		if age.(int64) != 31 {
			t.Fatalf("age = %v, want 31", age)
		}
	})
}

func TestRemoveSwapWithLast(t *testing.T) {
	runWithAdapters(t, func(t *testing.T, d deps.Deps) {
		users := newUsers(t, d)
		mustCreate(t, users, "u1@x.com", "u1", 1)
		middle := mustCreate(t, users, "u2@x.com", "u2", 2)
		mustCreate(t, users, "u3@x.com", "u3", 3)

		if e := middle.Remove(); e.Msg != "" {
			t.Fatalf("Remove: %v", e)
		}
		all, err := users.ListAll()
		if err != nil {
			t.Fatalf("ListAll: %v", err)
		}
		if len(all) != 2 {
			t.Fatalf("size after delete = %d, want 2", len(all))
		}
		// Last record (u3) must have been swapped into position 1..2.
		emails := map[string]bool{}
		for _, u := range all {
			email, e := u.Get("email")
			if e != nil {
				t.Fatalf("Get(email): %v", e)
			}
			emails[email.(string)] = true
		}
		if !emails["u1@x.com"] || !emails["u3@x.com"] || emails["u2@x.com"] {
			t.Fatalf("unexpected survivors: %v", emails)
		}
		// Victim's index entries must be gone, so the email is reusable.
		if users.FindByKey("email", "u2@x.com") != nil {
			t.Fatal("deleted record still resolvable by email")
		}
		mustCreate(t, users, "u2@x.com", "u2", 2)
		// Double remove is a no-op.
		if e := middle.Remove(); e.Msg != "" {
			t.Fatalf("second Remove should be a no-op, got %v", e)
		}
		// Ids are never reused: 3 creations + 1 more = last id 4.
		fresh := mustCreate(t, users, "u5@x.com", "u5", 5)
		if fresh.Id() != 5 {
			t.Fatalf("new id = %d, want 5 (ids never reused)", fresh.Id())
		}
	})
}

func TestPagination(t *testing.T) {
	runWithAdapters(t, func(t *testing.T, d deps.Deps) {
		users := newUsers(t, d)
		for i := 1; i <= 5; i++ {
			mustCreate(t, users, string(rune('a'+i))+"@x.com", string(rune('a'+i)), i)
		}
		page, err := users.List(2, 2)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(page) != 2 {
			t.Fatalf("page size = %d, want 2", len(page))
		}
		// Past-the-end page is empty, not an error.
		page, err = users.List(6, 10)
		if err != nil || len(page) != 0 {
			t.Fatalf("past-the-end page = %v, %v; want empty", page, err)
		}
	})
}

func TestSubDatabase(t *testing.T) {
	runWithAdapters(t, func(t *testing.T, d deps.Deps) {
		users := newUsers(t, d)
		user := mustCreate(t, users, "a@x.com", "alice", 30)

		if _, e := user.NewSubItem("sessions", map[string]any{"token": "t1", "creation": 100}); e != nil {
			t.Fatalf("NewSubItem: %v", e)
		}
		if _, e := user.NewSubItem("sessions", map[string]any{"token": "t2", "creation": 200}); e != nil {
			t.Fatalf("NewSubItem: %v", e)
		}
		sessions := user.ListAll("sessions")
		if len(sessions) != 2 {
			t.Fatalf("sessions = %d, want 2", len(sessions))
		}
		token, e := sessions[0].Get("token")
		if e != nil || token.(string) != "t1" {
			t.Fatalf("token = %v, %v; want t1", token, e)
		}
		// Duplicate session token conflicts.
		if _, e := user.NewSubItem("sessions", map[string]any{"token": "t1", "creation": 300}); e == nil {
			t.Fatal("duplicate sub token should conflict")
		}
		// Removing the parent clears the sub-database too.
		if e := user.Remove(); e.Msg != "" {
			t.Fatalf("Remove: %v", e)
		}
		recreated := mustCreate(t, users, "a@x.com", "alice", 30)
		if got := recreated.ListAll("sessions"); len(got) != 0 {
			t.Fatalf("sessions of recreated user = %d, want 0", len(got))
		}
	})
}
