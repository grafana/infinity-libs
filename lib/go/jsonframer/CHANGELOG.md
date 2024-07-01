# @grafana/infinity-jsonframer

## 1.1.0

- removed sqlite3 framer

## 1.0.0

- chore release

## 0.2.3

- Fixed a bug in multi framer where long frame conversion

## 0.2.2

- Support for timeseries in multi framer

## 0.2.1

- Fixed a bug in JSON multi framer

## 0.2.0

- Support for multi frames. New function `func ToFrames(jsonString string, options FramerOptions) (frames []*data.Frame, err error)` added to return multi frame response.

## 0.1.1

- 873e734: cleanup

## 0.1.0

- Allow column overrides

## 0.0.5

- replaced the backend package `blues/jsonata-go` to `xiatechs/jsonata-go`

## 0.0.4

- 🐛 **Chore**: updated build dependency turbo to 1.10.6

## 0.0.3

- 🐛 **Chore**: Fixed an issue with the github actions
