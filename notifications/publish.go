package notifications

// PublishMessage is just a wrapper around SNS adapter's publishing method
func (s *Service) PublishMessage(endpointARN, msg string, opt map[string]interface{}) (string, error) {
	return s.snsAdapter.PublishMessage(endpointARN, msg, opt)
}
