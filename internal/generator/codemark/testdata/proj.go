package codemark

// +codemark:option:summary="some summary"
// +codemark:option:description="some description"
type Opt1 string

// +codemark:option:summary="some summary"
// +codemark:option:description=`
// Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin in fermentum lacus, vitae blandit sem. Cras eget sem
// lorem. Fusce volutpat turpis in pellentesque tristique. Nullam varius eget sem malesuada dapibus. Vestibulum pretium accumsan mauris eu vulputate.
// Pellentesque tempor ultrices lacus, vel hendrerit nulla laoreet sit amet. Ut at sollicitudin nunc. Cras vel tortor at felis aliquet mattis nec vel
// lacus.`
type Opt2 string
