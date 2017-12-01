package externalmocks

// from vendor
//go:generate mockery -output . -outpkg externalmocks -case underscore -dir ../../vendor/gitlab.com/swarmfund/go/doorman -name Doorman
//go:generate mockery -output . -outpkg externalmocks -case underscore -dir ../../vendor/gitlab.com/swarmfund/go/doorman/data -name AccountQ

// from gopath
///go:generate mockery -output . -outpkg externalmocks -case underscore -dir $GOPATH/src/gitlab.com/swarmfund/go/doorman -name Doorman
///go:generate mockery -output . -outpkg externalmocks -case underscore -dir $GOPATH/src/gitlab.com/swarmfund/go/doorman/data -name AccountQ
