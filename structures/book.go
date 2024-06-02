package structures

type Book struct {
	ISBN            string `json:"isbn"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	Year            int    `json:"year"`
	Pages           int    `json:"pages"`
	Publisher       string `json:"publisher"`
	Language        string `json:"language"`
	Genre           string `json:"genre"`
	CopiesAvailable int    `json:"copies_available"`
}
