package message

import (
	"os"
	"reflect"
	"testing"
)

var msgServer *MessageServer

func TestMessageServer_Add(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want MessageObj
	}{
		{"positive flow 1",
			args{"first"},
			MessageObj{1, "first", nil},
		},
		{"negative flow",
			args{"second"},
			MessageObj{1, "second", nil}, // id should be 2
		},
		{"positive flow 2",
			args{"third"},
			MessageObj{3, "third", nil},
		},
	}
	for id, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if id == 1 {
				if got := msgServer.Add(tt.args.msg); reflect.DeepEqual(got, tt.want) {
					t.Errorf("Add() = %v, want %v", got, tt.want)
				}
			} else {
				if got := msgServer.Add(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Add() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMessageServer_GetAll_Pos(t *testing.T) {
	_, err := msgServer.GetAll()
	if err != nil {
		t.Error("GetAll() positive flow failed with error ,", err)
		return
	}
	if len(msgServer.msgStore) != 3 {
		t.Errorf("GetAll() positive flow failed as message store length %d is incorrect", len(msgServer.msgStore))
	}
}

func TestMessageServer_Get(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    MessageObj
		wantErr bool
	}{
		{"positive flow 1",
			args{1},
			MessageObj{1, "first", nil},
			false,
		},
		{"negative flow",
			args{5},
			MessageObj{},
			true,
		},
	}
	for id, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := msgServer.Get(tt.args.id)
			if id == 0 {
				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			} else {
				if err == nil {
					t.Error("Get() failed. Expected error to be non-nil")
					return
				}
			}
		})
	}
}

func TestMessageServer_Delete(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"positive flow 1",
			args{1},
			false,
		},
		{"negative flow",
			args{5},
			true,
		},
		{"positive flow 2",
			args{2},
			false,
		},
		{"positive flow 3",
			args{3},
			false,
		},
	}
	for id, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if id != 1 {
				if err := msgServer.Delete(tt.args.id); (err != nil) != tt.wantErr {
					t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err := msgServer.Delete(tt.args.id); err == nil {
					t.Errorf("Delete() failed. Expected error to be non-nil")
				}
			}
		})
	}
}

func TestMessageServer_GetAll_Neg(t *testing.T) {
	_, err := msgServer.GetAll()
	if err == nil {
		t.Error("GetAll() negative flow failed, no error returned")
		return
	}
	if len(msgServer.msgStore) > 0 {
		t.Error("GetAll() negative flow failed, length of message store non-zero")
		return
	}
}

func TestMain(m *testing.M) {
	msgServer = NewMessageServer()
	exitVal := m.Run()
	os.Exit(exitVal)
}
