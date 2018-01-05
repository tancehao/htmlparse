package htmlarse

import (
    "fmt"
)

func (t *Tag)Write(position int64, data []byte) (n, error) {}

func (t *Tag)WriteAfter(data []byte) (n, error) {}

func (t *Tag)WriteBefore(data []byte) (n, error) {}

func (t *Text)Write(position int64, data []byte) (n, error) {}
