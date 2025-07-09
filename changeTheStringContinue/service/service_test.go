package service

import (
	"errors"
	pres "github.com/MaksimHank/WorkingWithFile/changeTheStringContinue/presenter"
	prod "github.com/MaksimHank/WorkingWithFile/changeTheStringContinue/producer"
	"github.com/MaksimHank/WorkingWithFile/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestService_changeTheStringToAsterisks(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "one link",
			input: "Hello, it's my website http://example.com",
			want:  "Hello, it's my website http://***********",
		},
		{
			name:  "multiple links",
			input: "Hello, it's my website http://example.com; And this too http://test.com; http://123",
			want:  "Hello, it's my website http://***********; And this too http://********; http://***",
		},
		{
			name:  "without links",
			input: "Just regular text",
			want:  "Just regular text",
		},
		{
			name:  "link with path",
			input: "Check http://domain.com/path?query=value",
			want:  "Check http://***************************",
		},
		{
			name:  "link at start",
			input: "http://start.com is first",
			want:  "http://********* is first",
		},
		{
			name:  "email should be ignored",
			input: "Contact email@domain.com",
			want:  "Contact email@domain.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			if got := s.changeTheStringToAsterisks(tt.input); got != tt.want {
				t.Errorf("changeTheStringToAsterisks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Run(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(prod *mocks.Producer, pres *mocks.Presenter)
		wantErr    bool
	}{
		{
			name: "success",
			setupMocks: func(prod *mocks.Producer, pres *mocks.Presenter) {
				prod.On("Produce").Return([]string{"Hello, it's my website http://example.com"}, nil)
				pres.On("Present", []string{"Hello, it's my website http://***********"}).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Produce Error",
			setupMocks: func(prod *mocks.Producer, pres *mocks.Presenter) {
				prod.On("Produce").Return(nil, errors.New("produce failed"))
			},
			wantErr: true,
		},
		{
			name: "Present error",
			setupMocks: func(prod *mocks.Producer, pres *mocks.Presenter) {
				prod.On("Produce").Return([]string{"Hello, it's my website http://***********"}, nil)
				pres.On("Present", []string{"Hello, it's my website http://***********"}).
					Return(errors.New("present failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prodMock := new(mocks.Producer)
			presMock := new(mocks.Presenter)
			tt.setupMocks(prodMock, presMock)
			s := &Service{
				prod: prodMock,
				pres: presMock,
			}
			if err := s.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			prodMock.AssertExpectations(t)
			presMock.AssertExpectations(t)
		})
	}
}

func TestFileProducer_Produce(t *testing.T) {
	createTempFile := func(content string) string {
		tmpFile, err := os.CreateTemp("", "testfile")
		require.NoError(t, err)
		_, err = tmpFile.WriteString(content)
		require.NoError(t, err)
		tmpFile.Close()
		return tmpFile.Name()
	}

	tests := []struct {
		name        string
		inputFile   string
		fileContent string
		want        []string
		wantErr     string
	}{
		{
			name:        "Successful processing",
			inputFile:   createTempFile("Hello http://example.com"),
			fileContent: "Hello http://example.com",
			want:        []string{"Hello http://example.com"},
		},
		{
			name:        "Multiple lines",
			inputFile:   createTempFile("Line1\nLine2 http://test.com"),
			fileContent: "Line1\nLine2 http://test.com",
			want:        []string{"Line1", "Line2 http://test.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := &prod.FileProducer{
				InputFile: tt.inputFile,
			}
			got, err := fp.Produce()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFilePresenter_Present(t *testing.T) {
	readFile := func(path string) string {
		data, err := os.ReadFile(path)
		if err != nil {
			return ""
		}
		return string(data)
	}

	tests := []struct {
		name        string
		outputFile  string
		data        []string
		wantContent string
		wantErr     bool
	}{
		{
			name:        "Successful write",
			outputFile:  "",
			data:        []string{"Hello, it's my website http://***********;", "And this too http://********"},
			wantContent: "Hello, it's my website http://***********;\nAnd this too http://********\n",
			wantErr:     false,
		},
		{
			name:       "Write to invalid path",
			outputFile: "/invalid_directory/output.txt",
			data:       []string{"Hello, it's my website http://***********"},
			wantErr:    true,
		},
		{
			name:        "Empty data",
			outputFile:  "",
			data:        []string{},
			wantContent: "",
			wantErr:     false,
		},
		{
			name:        "Special characters",
			outputFile:  "",
			data:        []string{"line with \n newline", "Contact email@domain.com"},
			wantContent: "line with \n newline\nContact email@domain.com\n",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputPath := tt.outputFile
			if outputPath == "" && !tt.wantErr {
				tmpFile, err := os.CreateTemp("", "test_output")
				require.NoError(t, err)
				tmpFile.Close()
				outputPath = tmpFile.Name()
				defer os.Remove(outputPath)
			}

			fp := &pres.FilePresenter{
				OutputFile: outputPath,
			}

			err := fp.Present(tt.data)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tt.wantContent != "" {
				content := readFile(outputPath)
				assert.Equal(t, tt.wantContent, content)
			}
		})
	}
}
