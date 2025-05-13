package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	currentDirectory   = ""
	binaryName         = "mado"
	binaryPath         = ""
	inputDirectory     = "input"
	expectedDirectory  = "expected"
	actualDirectory    = "actual"
	auxiliaryDirectory = "auxiliary"
)

type UnnamedTest struct {
	args     []string
	expected string
}

type Test struct {
	name     string
	args     []string
	expected string
}

func TestMain(m *testing.M) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	currentDirectory = dir

	err = os.Chdir("..")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	dir, err = os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	binaryPath = filepath.Join(dir, binaryName)
	os.Exit(m.Run())
}

func runBinary(args []string) ([]byte, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")
	return cmd.CombinedOutput()
}

func getExpectedOutput(fixture string) ([]byte, error) {
	path := filepath.Join(currentDirectory, expectedDirectory, fixture)
	return os.ReadFile(path)
}

func getActualOutput(fixture string) ([]byte, error) {
	path := filepath.Join(currentDirectory, actualDirectory, fixture)
	return os.ReadFile(path)
}

func asInput(file string) string {
	return filepath.Join(currentDirectory, inputDirectory, file)
}

func asAuxiliary(file string) string {
	return filepath.Join(currentDirectory, auxiliaryDirectory, file)
}

func TestCliArgs(t *testing.T) {
	tests :=
		[]Test{
			{"replace simple", []string{"replace", "-i", asInput("md1.md"), "-r", asAuxiliary("replacements.txt")}, "replace.md"},
			{"html 1", []string{"html", "-i", asInput("md1.md")}, "html1.html"},
			{"html 2", []string{"html", "-i", asInput("md2.md")}, "html2.html"},
			{"html replace", []string{"html", "-i", asInput("md1.md"), "-r", asAuxiliary("replacements.txt")}, "html-replace.html"},
			{"jira 1", []string{"jira", "-i", asInput("md1.md")}, "jira1"},
			{"jira 2", []string{"jira", "-i", asInput("md2.md")}, "jira2"},
			{"jira replace", []string{"jira", "-i", asInput("md1.md"), "-r", asAuxiliary("replacements.txt")}, "jira-replace"},
		}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args = append(tt.args, "-o", filepath.Join(currentDirectory, actualDirectory, tt.expected))
			_, err := runBinary(tt.args)
			if err != nil {
				t.Fatal(err)
			}

			actual, err := getActualOutput(tt.expected)
			if err != nil {
				t.Fatal(err)
			}

			expected, err := getExpectedOutput(tt.expected)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("actual = %s, expected = %s", actual, expected)
			}
		})
	}
}

func TestStdout(t *testing.T) {
	test := Test{"replace simple", []string{"replace", "-i", asInput("md1.md"), "-r", asAuxiliary("replacements.txt")}, "replace.md"}
	actual, err := runBinary(test.args)
	if err != nil {
		t.Fatal(err)
	}

	_, err = getActualOutput(test.expected)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := getExpectedOutput(test.expected)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual = %s, expected = %s", actual, expected)
	}
}

func TestForce(t *testing.T) {
	test := UnnamedTest{[]string{"replace", "-i", asInput("force.md"), "-r", asAuxiliary("replacements.txt"), "-f"}, "force.md"}
	test.args = append(test.args, "-o", filepath.Join(currentDirectory, actualDirectory, test.expected))

	for range 2 {
		_, err := runBinary(test.args)
		if err != nil {
			t.Fatal(err)
		}

		actual, err := getActualOutput(test.expected)
		if err != nil {
			t.Fatal(err)
		}

		expected, err := getExpectedOutput(test.expected)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("actual = %s, expected = %s", actual, expected)
		}
	}
}

func TestConfig(t *testing.T) {
	test := UnnamedTest{[]string{"replace", "-i", asInput("md1.md"), "-c", asAuxiliary("replace.env")}, "replace-config.md"}
	test.args = append(test.args, "-o", filepath.Join(currentDirectory, actualDirectory, test.expected))

	_, err := runBinary(test.args)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := getActualOutput(test.expected)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := getExpectedOutput(test.expected)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual = %s, expected = %s", actual, expected)
	}
}

func TestEnv(t *testing.T) {
	test := UnnamedTest{[]string{"replace", "-i", asInput("md1.md")}, "replace-env.md"}
	test.args = append(test.args, "-o", filepath.Join(currentDirectory, actualDirectory, test.expected))

	cmd := exec.Command(binaryPath, test.args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata", "MADO_REPLACE="+asAuxiliary("replacements.txt"))
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	actual, err := getActualOutput(test.expected)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := getExpectedOutput(test.expected)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual = %s, expected = %s", actual, expected)
	}
}
