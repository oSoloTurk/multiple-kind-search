package main

type Author struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Bio       string `json:"bio,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type News struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	AuthorID  string   `json:"author_id"`
	Tags      []string `json:"tags,omitempty"`
	ImageURL  string   `json:"image_url,omitempty"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}
