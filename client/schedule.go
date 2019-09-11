package client

type Schedule struct {
	Period          string `yaml:"period"`
	Mon             bool   `yaml:"mon"`
	Tues            bool   `yaml:"tues"`
	Wed             bool   `yaml:"wed"`
	Thurs           bool   `yaml:"thurs"`
	Fri             bool   `yaml:"fri"`
	Sat             bool   `yaml:"sat"`
	Sun             bool   `yaml:"sun"`
	Time            string `yaml:"time"`
	TimeZone        string `yaml:"timeZone"`
	ExcludeHolidays string `yaml:"excludeHolidays,omitempty"`
}
