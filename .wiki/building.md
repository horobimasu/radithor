# notes
1. make sure you have [bun](https://bun.sh/) (at least version 1.4.2) installed
2. make sure you have [golang](https://go.dev/dl/) (at least version 1.26.3) installed
3. you can technically install whatever js runtime you want, just change the build script where bun is used

# building live
0. make sure you have bun and golang installed as stated above
1. run the build script `/scripts/build-live.js`
2. a folder called `.live` should appear in the root directory with some stuff inside
3. run `builder.exe` and go to the url it shows
4. if you want, you can delete all other folders and files, you will only need the live folder from now on (you are gonna have to redownload everything for updates anyways) (dont delete the radithor folder tho, a `.builds` folder will be created there)

# building radithor
0. you should be on a localhost website after clicking that url from `builder.exe`
1. fill in the data with your config
2. click the "build" button on the bottom left
3. a folder called `.builds` should appear in the root directory
4. in there will be the radithor build, the latest build will always be prefixed with "latest-"

# faq
1. what is the website code so shit? it was not originally made to be bundled into a single file so i just mashed everything up into a export default thing to get it over with, i plan to rewrite it some day
