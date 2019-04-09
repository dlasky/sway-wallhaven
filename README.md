# sway-wallhaven

### Basic usage

`go build .`

`./wallhaven fetch`

`./wallhaven set`

### Using in sway

in your sway config instead of:

`output "*" background image.jpg fill`

try

`exec wallhaven restore`

then its possible to bind fetch and set to unused keys

