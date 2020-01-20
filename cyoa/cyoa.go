package github.com/MrBTTF/gophercises/cyoa

type Option struct {
	Text    string
	NextArc string `json:"arc"`
}

type Arc struct {
	Title   string
	Story   []string
	Options []Option
}

func LoadBook(filename string) (map[string]Arc, error) {
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var arcs map[string]Arc
	json.Unmarshal(jsonData, &arcs)

	return arcs, nil
}