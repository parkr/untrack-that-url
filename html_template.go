package untrackthaturl

import "io"

func RenderHTML(w io.Writer) error {
	_, err := w.Write([]byte(`not implemented yet`))
	return err
}
