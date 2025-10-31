#import "../lib/libDocumentiGenerici.typ": *

#let meta(
  versione: none,
  htmlId: none,
  titolo: none,
  body
) = {
  set document(
    title: titolo,
    author: "",
    keywords: (htmlId, versione),
  )
  
  body
}