package pac

type producer struct {
	update string
	app    Application
}

func (p producer) Produce(data interface{}) error {
	sparql, err := p.app.profile.GetUpdate(p.update, data)

	if err != nil {
		return err
	}

	return p.app.client.Update(sparql)
}
