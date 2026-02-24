package cli

import (
	"bytes"
	"fmt"
	"testing"
	"text/tabwriter"
	"time"

	"github.com/EshkinKot1980/GophKeeper/internal/client/cli/mocks"
	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_list(t *testing.T) {
	type listItem struct {
		c1 uint64
		c2 string
		c3 string
		c4 string
		c5 time.Time
	}

	now := time.Now()
	successList := []listItem{
		{c1: 10, c2: dto.SecretTypeText, c3: "name10", c5: now.Add(-24 * time.Hour)},
		{c1: 13, c2: dto.SecretTypeFile, c3: "name13", c4: "secret.txt", c5: now},
	}

	successOut := new(bytes.Buffer)
	w := tabwriter.NewWriter(successOut, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tType\tNname\tFileName\tCreated")
	for _, item := range successList {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			item.c1,
			item.c2,
			item.c3,
			item.c4,
			item.c5.Format("2006-01-02 15:04:05"),
		)
	}
	w.Flush()

	type want struct {
		output string
		err    string
	}

	tests := []struct {
		name  string
		setup func(t *testing.T) SecretService
		want  want
	}{
		{
			name: "success",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().InfoList().
					Return(
						[]dto.SecretInfo{
							{
								ID:       10,
								DataType: dto.SecretTypeText,
								Name:     "name10",
								Meta:     []dto.MetaData{},
								Created:  now.Add(-24 * time.Hour),
							},
							{
								ID:       13,
								DataType: dto.SecretTypeFile,
								Name:     "name13",
								Meta: []dto.MetaData{
									{Name: "some meta", Value: "some value"},
									{Name: MetaFileName, Value: "secret.txt"},
								},
								Created: now,
							},
						},
						nil,
					)
				return service
			},
			want: want{
				output: successOut.String(),
			},
		},
		{
			name: "get_list_error",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().InfoList().
					Return(nil, fmt.Errorf("get list error"))
				return service
			},
			want: want{
				err: "get list error",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			secretService = test.setup(t)

			out := new(bytes.Buffer)
			err := list(out)

			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, test.want.err, gotErr, "List error")
			assert.Equal(t, test.want.output, out.String(), "List output")
		})
	}
}
