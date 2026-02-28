package testdata

type AutoIncModel struct {
	Id       int64  `db:"id,pk"`
	SeqNo    int64  `db:"seq_no,autoincrement"`
	Name     string `db:"name"`
	Birthday string `db:"birthday"`
}
