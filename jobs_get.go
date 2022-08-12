package gojobs

func (j *JobsGorm) getCornKeyIp() string {
	return j.config.cornPrefix + "_" + j.config.outsideIp
}

func (j *JobsGorm) getCornKeyChannel() string {
	return j.config.cornKeyIp
}

func (j *JobsGorm) getCornKeyChannels() string {
	return j.config.cornKeyIp + "_*"
}
