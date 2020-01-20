package scopr

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type User struct {
	Username    string   `json:"username" scope:"user,admin,root"`
	Password    string   `json:"-"`
	EmptyField  string   `json:"empty"`
	Email       string   `json:"email" scope:"user,admin,root"`
	Secrets     string   `json:"secrets" scope:"admin,root"`
	SuperSecret string   `json:"super_secret" scope:"root"`
	Admin       NullBool `json:"admin" scope:"user"`
}

var exampleUser = User{
	Username:    "hunterlong",
	Password:    "password12345",
	EmptyField:  "yoyoyo",
	Email:       "info@wowowww.com",
	Secrets:     "$03043iod2djswuhwie",
	SuperSecret: "this is absolutly top secret",
	Admin:       NullBool{Bool: true, Valid: true},
}

func sliceUsers(amount int) []*User {
	var users []*User
	for i := 0; i < amount; i++ {
		users = append(users, &exampleUser)
	}
	return users
}

type NullBool sql.NullBool

func (v NullBool) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Bool)
	} else {
		return json.Marshal(false)
	}
}

func (v *NullBool) UnmarshalJSON(data []byte) error {
	var x *bool
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Bool = *x
	} else {
		v.Valid = false
	}
	return nil
}

func TestNativeScope(t *testing.T) {
	dd := New(exampleUser, "admin")
	data, err := json.Marshal(dd)
	require.Nil(t, err)
	t.Log(string(data))
	var scopedJson User
	json.Unmarshal(data, &scopedJson)
	assert.Equal(t, exampleUser.Email, scopedJson.Email)
	assert.Equal(t, "", scopedJson.SuperSecret)
}

func TestRootScope(t *testing.T) {
	scopedJson, err := Json(exampleUser, "root")
	require.Nil(t, err)
	assert.Equal(t, exampleUser.Email, toUser(scopedJson).Email)
	assert.Equal(t, "this is absolutly top secret", toUser(scopedJson).SuperSecret)
}

func TestAdminScope(t *testing.T) {
	scopedJson, err := Json(exampleUser, "admin")
	require.Nil(t, err)
	assert.Equal(t, exampleUser.Email, toUser(scopedJson).Email)
	assert.Equal(t, "", toUser(scopedJson).SuperSecret)
}

func TestUserScope(t *testing.T) {
	scopedJson, err := Json(exampleUser, "user")
	require.Nil(t, err)
	assert.Equal(t, exampleUser.Email, toUser(scopedJson).Email)
	assert.Equal(t, "", toUser(scopedJson).SuperSecret)
	assert.Equal(t, true, toUser(scopedJson).Admin.Bool)
	t.Log(string(scopedJson))
}

func TestEmptyScope(t *testing.T) {
	scopedJson, err := Json(exampleUser, "")
	require.Nil(t, err)
	assert.Equal(t, "", toUser(scopedJson).Email)
	assert.Equal(t, "", toUser(scopedJson).SuperSecret)
}

func TestUserSliceScope(t *testing.T) {
	users := sliceUsers(10)
	dd := New(users, "admin")
	data, err := json.Marshal(dd)
	require.Nil(t, err)

	var allUsers []*User
	err = json.Unmarshal(data, &allUsers)
	require.Nil(t, err)

	assert.Equal(t, 10, len(allUsers))
}

func TestEncodeWriter(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	err := enc.Encode(exampleUser, "admin")
	require.Nil(t, err)

	user := toUser(buf.Bytes())
	assert.Equal(t, exampleUser.Email, user.Email)
	assert.Equal(t, "", user.SuperSecret)
}

func toUser(d []byte) User {
	var u User
	json.Unmarshal(d, &u)
	return u
}
