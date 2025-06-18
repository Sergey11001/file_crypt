package errs

import (
	"errors"
	"testing"
)

func init() {
	Secure = true
}

func TestClassIs(t *testing.T) {
	type args struct {
		class Class
		err   error
	}
	type want struct {
		result bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "regular",
			args: args{
				class: Invalid,
				err:   Invalid.New("Error", "error"),
			},
			want: want{
				result: true,
			},
		},
		{
			name: "native",
			args: args{
				class: Invalid,
				err:   errors.New("error"),
			},
			want: want{
				result: false,
			},
		},
		{
			name: "unclassified",
			args: args{
				class: Unclassified,
				err:   Unclassified.New("Error", "error"),
			},
			want: want{
				result: false,
			},
		},
		{
			name: "native_internal",
			args: args{
				class: Internal,
				err:   errors.New("error"),
			},
			want: want{
				result: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClassIs(tt.args.class, tt.args.err)
			if result != tt.want.result {
				t.Errorf("%t expected, got %t", tt.want.result, result)
			}
		})
	}
}

func TestCodeIs(t *testing.T) {
	type args struct {
		code string
		err  error
	}
	type want struct {
		result bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "regular",
			args: args{
				code: "Error",
				err:  Invalid.New("Error", "error"),
			},
			want: want{
				result: true,
			},
		},
		{
			name: "native",
			args: args{
				code: "Error",
				err:  errors.New("error"),
			},
			want: want{
				result: false,
			},
		},
		{
			name: "empty_code",
			args: args{
				code: "",
				err:  Invalid.New("Error", "error"),
			},
			want: want{
				result: false,
			},
		},
		{
			name: "unknown",
			args: args{
				code: Unknown,
				err:  Invalid.New(Unknown, "error"),
			},
			want: want{
				result: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CodeIs(tt.args.code, tt.args.err)
			if result != tt.want.result {
				t.Errorf("%t expected, got %t", tt.want.result, result)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		class   Class
		code    string
		message string
		details map[string]string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "regular",
			args: args{
				err: Invalid.New("Error", "error", Details{
					"key1": "value1",
					"key2": "value2",
				}),
			},
			want: want{
				class:   Invalid,
				code:    "Error",
				message: "error",
				details: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
		{
			name: "native",
			args: args{
				err: errors.New("error"),
			},
			want: want{
				class:   Internal,
				code:    Unknown,
				message: "unknown",
				details: nil,
			},
		},
		{
			name: "class_only",
			args: args{
				err: &errorWithClass{errors.New("error"), Invalid},
			},
			want: want{
				class:   Invalid,
				code:    Unknown,
				message: "unknown",
				details: nil,
			},
		},
		{
			name: "code_only",
			args: args{
				err: &errorWithCode{errors.New("error"), "Error"},
			},
			want: want{
				class:   Internal,
				code:    "Error",
				message: "error",
				details: nil,
			},
		},
		{
			name: "empty_code_only",
			args: args{
				err: &errorWithCode{errors.New("error"), ""},
			},
			want: want{
				class:   Internal,
				code:    Unknown,
				message: "unknown",
				details: nil,
			},
		},
		{
			name: "details_only",
			args: args{
				err: &errorWithDetails{errors.New("error"), map[string]string{
					"key1": "value1",
					"key2": "value2",
				}},
			},
			want: want{
				class:   Internal,
				code:    Unknown,
				message: "unknown",
				details: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
		{
			name: "empty_details_only",
			args: args{
				err: &errorWithDetails{errors.New("error"), nil},
			},
			want: want{
				class:   Internal,
				code:    Unknown,
				message: "unknown",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			class, code, message, details := Parse(tt.args.err)
			if class != tt.want.class {
				t.Errorf("class: %s expected, got %s", tt.want.class, class)
			}
			if code != tt.want.code {
				t.Errorf("code: %s expected, got %s", tt.want.code, code)
			}
			if message != tt.want.message {
				t.Errorf("message: %s expected, got %s", tt.want.message, message)
			}
			if !mapsEqual(details, tt.want.details) {
				t.Errorf("details: %v expected, got %v", tt.want.details, details)
			}
		})
	}
}

func TestNested(t *testing.T) {
	class1, code1, message1, detailers1, details1 := Class("Class1"), "Error1", "error1", []Detailer{
		DetailerFunc(func(details map[string]string) {
			details["key1"] = "value1"
		}),
	}, map[string]string{
		"key1": "value1",
	}
	class2, code2, message2, detailers2, details2 := Class("Class2"), "Error2", message1, []Detailer{
		Details{
			"key2": "value2",
			"key3": "value3",
		},
	}, map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	class3, code3, message3, detailers3, details3 := Class("Class3"), code2, message2, []Detailer{
		DetailerFunc(func(details map[string]string) {
			delete(details, "key1")
		}),
	}, map[string]string{
		"key2": "value2",
		"key3": "value3",
	}
	class4, code4, message4, detailers4, details4 := class3, code3, message3, []Detailer{
		(Details)(nil),
	}, details3
	class5, code5, message5, detailers5, details5 := class4, code4, message4, []Detailer{
		DetailerFunc(func(details map[string]string) {
			clear(details)
		}),
	}, (map[string]string)(nil)

	err1 := class1.New(code1, message1, detailers1...)
	err2 := class2.As(code2, err1, detailers2...)
	err3 := class3.Cast(err2, detailers3...)
	err4 := Detail(err3, detailers4...)
	err5 := Detail(err4, detailers5...)

	class, code, message, details := Parse(err1)
	if class != class1 {
		t.Errorf("err1 class: %s expected, got %s", class1, class)
	}
	if code != code1 {
		t.Errorf("err1 code: %s expected, got %s", code1, code)
	}
	if message != message1 {
		t.Errorf("err1 message: %q expected, got %q", message1, message)
	}
	if !mapsEqual(details, details1) {
		t.Errorf("err1 details: %v expected, got %v", details1, details)
	}

	class, code, message, details = Parse(err2)
	if class != class2 {
		t.Errorf("err2 class: %s expected, got %s", class2, class)
	}
	if code != code2 {
		t.Errorf("err2 code: %s expected, got %s", code2, code)
	}
	if message != message2 {
		t.Errorf("err2 message: %q expected, got %q", message2, message)
	}
	if !mapsEqual(details, details2) {
		t.Errorf("err2 details: %v expected, got %v", details2, details)
	}

	class, code, message, details = Parse(err3)
	if class != class3 {
		t.Errorf("err3 class: %s expected, got %s", class3, class)
	}
	if code != code3 {
		t.Errorf("err3 code: %s expected, got %s", code3, code)
	}
	if message != message3 {
		t.Errorf("err3 message: %q expected, got %q", message3, message)
	}
	if !mapsEqual(details, details3) {
		t.Errorf("err3 details: %v expected, got %v", details3, details)
	}

	class, code, message, details = Parse(err4)
	if class != class4 {
		t.Errorf("err4 class: %s expected, got %s", class4, class)
	}
	if code != code4 {
		t.Errorf("err4 code: %s expected, got %s", code4, code)
	}
	if message != message4 {
		t.Errorf("err4 message: %q expected, got %q", message4, message)
	}
	if !mapsEqual(details, details4) {
		t.Errorf("err4 details: %v expected, got %v", details4, details)
	}

	class, code, message, details = Parse(err5)
	if class != class5 {
		t.Errorf("err5 class: %s expected, got %s", class5, class)
	}
	if code != code5 {
		t.Errorf("err5 code: %s expected, got %s", code5, code)
	}
	if message != message5 {
		t.Errorf("err5 message: %q expected, got %q", message5, message)
	}
	if !mapsEqual(details, details5) {
		t.Errorf("err5 details: %v expected, got %v", details5, details)
	}
}

func mapsEqual[M1, M2 ~map[K]V, K, V comparable](m1 M1, m2 M2) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}

	return true
}
