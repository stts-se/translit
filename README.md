# translit

A collection of utilities for transliterating non-Latin scripts into Latin (and reverse, where applicable).

They are not synched in default behaviour, but that might happen later on.

Work in progress.


## Installation

1. Set up `go`

     Download: https://golang.org/dl/ (1.13 or higher)   
     Installation instructions: https://golang.org/doc/install             

2. Clone the source code

   `$ git clone https://github.com/stts-se/translit.git`  
   `$ cd translit`   
   
3. Test (optional)

   `translit$ go test ./...`


4. Pre-compile binaries

    `translit$ go install ./...`


---

## Language versions

### Arabic Buckwalter

 `translit$ buckwalter <arabic text>`

References:
  * http://www.qamus.org/transliteration.htm
  * https://en.wikipedia.org/wiki/Buckwalter_transliteration

### Farsi

EI (2012)

 `translit$ far2lat <farsi text>`

References:
  * https://en.wikipedia.org/wiki/Romanization_of_Persian

### Greek to Latin

Simplified version of ALA-LC [3]

 `translit$ grc2lat <greek text>`


References:
   * https://en.wikipedia.org/wiki/Romanization_of_Greek#Modern_Greek


### Russian to Latin

A simplified version of the 'Road signs' system.

For Swedish style transliteration, we are using a simplified version TT's recommendations (link below).

 `translit$ rus2lat <russian text>`


References:
* https://en.wikipedia.org/wiki/Romanization_of_Russian
* https://tt.se/tt-spraket/ord-och-begrepp/internationellt/andra-sprak/ryska/

### Tamil to Latin

ISO 15919

 `translit$ tamil2lat <tamil text>`

References:
* https://en.wikipedia.org/wiki/Tamil_script

