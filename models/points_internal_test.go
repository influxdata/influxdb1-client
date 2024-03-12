package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalPointNoFields(t *testing.T) {
	points, err := ParsePointsString("m,k=v f=0i")
	if err != nil {
		t.Fatal(err)
	}

	// It's unclear how this can ever happen, but we've observed points that were marshaled without any fields.
	points[0].(*point).fields = []byte{}

	if _, err := points[0].MarshalBinary(); err != ErrPointMustHaveAField {
		t.Fatalf("got error %v, exp %v", err, ErrPointMustHaveAField)
	}
}

func TestBinaryField(t *testing.T) {
	t.Run("parse-b64-binary-field", func(t *testing.T) {
		buf := []byte(`m1 f_d_1st="MTIzCg=="b,f_s="some-string",f_d_middle="MTIzCg=="b,f_b=F,f_f=1.0,f_i=2i,f_a=[1i,2i],f_u=32u,f_d_last="MTIzCg=="b 123`)
		pt, err := parsePoint(buf, time.Now(), "n")
		require.NoError(t, err)

		fields, err := pt.Fields()
		assert.NoError(t, err)

		assert.Equal(t, []byte("123\n"), fields["f_d_1st"])
		assert.Equal(t, []byte("123\n"), fields["f_d_middle"])
		assert.Equal(t, []byte("123\n"), fields["f_d_last"])

		assert.Equal(t, "some-string", fields["f_s"])
		assert.Equal(t, false, fields["f_b"])
		assert.Equal(t, 1.0, fields["f_f"])
		assert.Equal(t, int64(2), fields["f_i"])
		assert.Equal(t, []int64{1, 2}, fields["f_a"])
		assert.Equal(t, uint64(32), fields["f_u"])

		for k, v := range fields {
			t.Logf("%s: %v", k, v)
		}
	})

	t.Run("parse-invalid-b64-field", func(t *testing.T) {
		buf := []byte(`m1 f_d="MTIzCg=="x,f_s="some-string",f_b=F,f_f=1.0,f_i=2i,f_a=[1i,2i],f_u=32u 123`)
		_, err := parsePoint(buf, time.Now(), "n")
		require.Error(t, err)
		t.Logf("expeted error: %s", err)
	})

	t.Run("parse-invalid-b64-string", func(t *testing.T) {
		buf := []byte(`m1 f_d="invalid-base-64-string"b,f_s="some-string",f_b=F,f_f=1.0,f_i=2i,f_a=[1i,2i],f_u=32u 123`)
		pt, err := parsePoint(buf, time.Now(), "n")
		assert.NoError(t, err)

		_, err = pt.Fields() // cant retrieve bad-base64-encoded field
		assert.Error(t, err)
		t.Logf("expeted error: %s", err)
	})

	t.Run("parse-invalid-b64-string", func(t *testing.T) {
		buf := []byte(`m1 f_d=""b,f_s="some-string",f_b=F,f_f=1.0,f_i=2i,f_a=[1i,2i],f_u=32u 123`)
		pt, err := parsePoint(buf, time.Now(), "n")
		assert.NoError(t, err)

		_, err = pt.Fields() // cant retrieve bad-base64-encoded field
		assert.Error(t, err)
		t.Logf("expeted error: %s", err)
	})

	t.Run("new-b64-field", func(t *testing.T) {
		binData := []byte(`hello:😄`)
		pt, err := NewPoint("some", nil, map[string]any{"f_d": binData}, time.Unix(0, 123))
		assert.NoError(t, err)

		lp := pt.String()
		t.Logf("line-proto: %s", lp)

		decodePt, err := parsePoint([]byte(lp), time.Now(), "n")
		assert.NoError(t, err)

		fields, err := decodePt.Fields()
		assert.NoError(t, err)
		assert.Equal(t, binData, fields["f_d"])
	})

	t.Run("binary-array", func(t *testing.T) {
		buf := []byte(`m1 f_d_arr=["MTIzCg=="b,"MTIzCg=="b, ] 123`)
		pt, err := parsePoint(buf, time.Now(), "n")
		require.NoError(t, err)

		fields, err := pt.Fields()
		_ = fields
		assert.NoError(t, err)

		assert.Equal(t, [][]byte{
			[]byte("123\n"),
			[]byte("123\n"),
		}, fields["f_d_arr"])

		for k, v := range fields {
			t.Logf("%s: %q", k, v)
		}
	})
}
