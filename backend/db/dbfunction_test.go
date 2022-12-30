package db

import (
	"fmt"
	"testing"
)

func TestInsertData(t *testing.T) {
	//userId := 12
	//postId := 1
	var tests = [][]any{
		//test correct input
		{"users", "marikh125", "msdf@gmail.com", "1234", "12.12.2022"},
		//test missing data
		{"posts", "marikh1", "first post", "hello world!", "12.12.2022"},
		//test missing table
		{"postds", "marikh1", "first post", "hello world!", "12.12.2022"},
		//{"comments", userId, postId, "first comment", "12.12.2022"}, ////////////////adding later
		{"topics", "art"},
		//test duplicate name
		{"topics", "art"},
		{"topics", "music"},
	}

	for _, tt := range tests {
		// t.Run enables running "subtests", one for each
		// table entry. These are shown separately
		// when executing `go test -v`.
		testname := fmt.Sprintf("%s,%s", tt[0], tt[1])
		t.Run(testname, func(t *testing.T) {
			if len(tt) == 5 {
				err := InsertData(tt[0].(string), tt[1], tt[2], tt[3], tt[4])
				if err != nil {
					t.Errorf("InsertData got error: %v:", err)

				}
			} else {
				err := InsertData(tt[0].(string), tt[1])
				if err != nil {
					t.Errorf("InsertData got error: %v:", err)

				}
			}
		})

	}

}

func TestSelectDataHandler(t *testing.T) {
	var tests = []struct {
		tableName, keyName, keyValue string

		want any
	}{
		//test correct input
		{"users", "userName", "marikh10", nil},
		{"users", "userName", "marikh6", nil},
		{"users12", "userName", "marikh3", nil},
		{"users", "email", "msdf@gmail.com", nil},
	}

	for _, tt := range tests {
		// t.Run enables running "subtests", one for each
		// table entry. These are shown separately
		// when executing `go test -v`.
		testname := fmt.Sprintf("%s,%s", tt.tableName, tt.keyValue)
		t.Run(testname, func(t *testing.T) {
			_, err := SelectDataHandler(tt.tableName, tt.keyName, tt.keyValue)
			if err != nil {
				t.Errorf("SelectDataHandler got error: %v:", err)

			}
			/* 	if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			} */
		})
	}

}
func TestDeleteData(t *testing.T) {
	var tests = []struct {
		tableName, key string
	}{
		//test correct input
		{"users", "marikh1"},
		//test missing data
		{"users", "marikh6"},
		//test missing table
		{"users121", "marikh3"},
	}
	for _, tt := range tests {
		// t.Run enables running "subtests", one for each
		// table entry. These are shown separately
		// when executing `go test -v`.
		testname := fmt.Sprintf("%s,%s", tt.tableName, tt.key)
		t.Run(testname, func(t *testing.T) {
			err := DeleteData(tt.tableName, tt.key)
			if err != nil {
				t.Errorf("DeleteData got error: %v:", err)

			}

		})
	}
}
