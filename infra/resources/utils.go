package resources

import "github.com/aws/jsii-runtime-go"

func sVtoP(e []string) *[]*string {
	var tmp []*string = []*string{}

	for _, v := range e {
		tmp = append(tmp, jsii.String(v))
	}

	return &tmp
}

func mVtoP(e map[string]string) *map[string]*string {
	var tmp = make(map[string]*string)

	for key, value := range e {
		tmp[key] = jsii.String(value)
	}

	return &tmp
}
