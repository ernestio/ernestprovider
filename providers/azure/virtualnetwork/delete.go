package virtualnetwork

// Delete : Deletes a nat object on azure
func (ev *Event) Delete() (err error) {
	resGroup := ev.ResourceGroupName
	name := ev.Name

	_, err = ev.client().Delete(resGroup, name, make(chan struct{}))

	return err
}
