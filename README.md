# go-img-converter
### Description
`go-img-converter` is open-source substitute for IMG2PNG software written in Go. It's used for converting images in IMG format, mainly used by NASA PDS to a more convenient format, like PNG (for now it only supports PNG). It can read 8-bit ant 16-bit IMG files
### Usage
`go-img-converter -input input.IMG -label label.LBL`
This command will convert `input.IMG` file based on `label.LBL` label and save as `input.png` in the same directory as original file.
### Building
`go build .`
### TODO list
Done? | Task
:---:| ---
⬜| Read and process .VIC files
⬜️| Automatically find label file for IMG files
⬜️| Multithreaded batch processing of directory
⬜️| Debayer color files