package scraper

//ProcessTitle helps in organising the task of getting the titles related to the given title
func ProcessTitle(title string) ([][]string, error) {
	op, err := mediaWikiCall(title)
	if err != nil {
		return nil, err
	}
	return getWikiTitles(op), nil

}
