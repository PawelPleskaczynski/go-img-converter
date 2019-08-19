# go-img-converter
### Description
`go-img-converter` is open-source substitute for IMG2PNG software written in Go. It's used for converting images in IMG format, mainly used by NASA PDS to a more convenient format, like PNG (for now it only supports PNG). It can read 8-bit and 16-bit IMG/VIC files (and 12-bit ones from InSight mission).

I tested conversion on images from following probes:
- Mars Science Laboratory
- InSight
- Juno
- Clementine 1
- I'll test with other sats/probes soon

### Usage
`go-img-converter -input input.IMG -label label.LBL`
This command will convert `input.IMG` file based on `label.LBL` label and save as `input.png` in the same directory as original file.
### Building
`go build .`
### TODO list
Done? | Task
:---:| ---
✅| Read and process .VIC files
✅| Automatically find image info in header
⬜| Remove header from file for processing 
⬜️| Automatically find label file for IMG files
⬜️| Multithreaded batch processing of directory
⬜️| Debayer color files