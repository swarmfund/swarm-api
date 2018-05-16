package externalmocks

// from vendor
//go:generate mockery -output . -outpkg externalmocks -case underscore -dir ../../vendor/gitlab.com/tokend/go/doorman -name Doorman
//go:generate mockery -output . -outpkg externalmocks -case underscore -dir ../../vendor/gitlab.com/tokend/go/doorman/data -name AccountQ

// from gopath
///go:generate mockery -output . -outpkg externalmocks -case underscore -dir $GOPATH/src/gitlab.com/tokend/go/doorman -name Doorman
///go:generate mockery -output . -outpkg externalmocks -case underscore -dir $GOPATH/src/gitlab.com/tokend/go/doorman/data -name AccountQ
