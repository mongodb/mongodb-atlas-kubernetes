package atlasproject

import (
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func Test_auditingInSync(t *testing.T) {
	type args struct {
		atlas *mongodbatlas.Auditing
		spec  *v1.Auditing
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Atlas and Operator Auditing are empty",
			args: args{
				atlas: nil,
				spec:  nil,
			},
			want: true,
		},
		{
			name: "Atlas Auditing is empty and Operator doesn't",
			args: args{
				atlas: nil,
				spec:  &v1.Auditing{Enabled: true},
			},
			want: false,
		},
		{
			name: "Operator Auditing is empty and Atlas doesn't",
			args: args{
				atlas: &mongodbatlas.Auditing{Enabled: toptr.MakePtr(true)},
				spec:  nil,
			},
			want: false,
		},
		{
			name: "Operator Auditing has different config from Atlas",
			args: args{
				atlas: &mongodbatlas.Auditing{
					AuditAuthorizationSuccess: toptr.MakePtr(false),
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					ConfigurationType:         "ReadOnly",
					Enabled:                   toptr.MakePtr(true),
				},
				spec: &v1.Auditing{
					AuditAuthorizationSuccess: true,
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					Enabled:                   true,
				},
			},
			want: false,
		},
		{
			name: "Operator Auditing has different config filter from Atlas",
			args: args{
				atlas: &mongodbatlas.Auditing{
					AuditAuthorizationSuccess: toptr.MakePtr(false),
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					ConfigurationType:         "ReadOnly",
					Enabled:                   toptr.MakePtr(true),
				},
				spec: &v1.Auditing{
					AuditAuthorizationSuccess: false,
					AuditFilter:               `{"atype":"authenticate","param":{"db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					Enabled:                   true,
				},
			},
			want: false,
		},
		{
			name: "Operator Auditing are Equal",
			args: args{
				atlas: &mongodbatlas.Auditing{
					AuditAuthorizationSuccess: toptr.MakePtr(false),
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					ConfigurationType:         "ReadOnly",
					Enabled:                   toptr.MakePtr(true),
				},
				spec: &v1.Auditing{
					AuditAuthorizationSuccess: false,
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					Enabled:                   true,
				},
			},
			want: true,
		},
		{
			name: "Operator Auditing are Equal when filter has newline in the end",
			args: args{
				atlas: &mongodbatlas.Auditing{
					AuditAuthorizationSuccess: toptr.MakePtr(false),
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					ConfigurationType:         "ReadOnly",
					Enabled:                   toptr.MakePtr(true),
				},
				spec: &v1.Auditing{
					AuditAuthorizationSuccess: false,
					AuditFilter:               "{\"atype\":\"authenticate\",\"param\":{\"user\":\"auditReadOnly\",\"db\":\"admin\",\"mechanism\":\"SCRAM-SHA-1\"}}\n",
					Enabled:                   true,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, auditingInSync(tt.args.atlas, tt.args.spec), "auditingInSync(%v, %v)", tt.args.atlas, tt.args.spec)
		})
	}
}
