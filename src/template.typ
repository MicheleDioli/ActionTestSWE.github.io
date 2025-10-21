#let primary = rgb("#2c3e50")
#let secondary = rgb("#000")
#let gray-dark = rgb("#2c3e50")

#let report(titolo: none, stato: none, versione: none, presenze: none, htmlId: none, body) = [

  #set document(
    title: [#titolo],
    author: "GlitchHub Team",
    keywords: htmlId,
  )

  #set page(
    paper: "a4",
    margin: (
      top: 3cm,
      bottom: 3cm,
      left: 2.5cm,
      right: 2.5cm,
    ),
  )

  #set text(
    font: ("Helvetica Neue", "Arial", "sans-serif"),
    size: 11pt,
    fill: gray-dark,
  )

  #set heading(numbering: "1.1")

  #set par(
    justify: true,
    leading: 0.75em,
    spacing: 1.2em,
  )


  #show heading.where(level: 1): it => {
    v(0.5em)
    rect(
      width: 100%,
      fill: primary,
      inset: 0pt,
      outset: 0pt,
    )[
      #block(inset: (x: 15pt, y: 12pt))[
        #text(size: 14pt, weight: "bold", fill: white)[
          #if it.numbering != none [
            #counter(heading).display(it.numbering)
            #h(0.5em)
          ]
          #it.body
        ]
      ]
    ]
    v(1em)
  }


  #show heading.where(level: 2): it => {
    v(1em)
    block[
      #text(size: 12pt, weight: "bold", fill: primary)[
        #if it.numbering != none [
          #counter(heading).display(it.numbering)
          #h(0.5em)
        ]
        #it.body
      ]
    ]
    v(0.5em)
  }


  #align(center)[
    #v(3cm)

    #align(center)[
      #image("./asset/lg-removebg-preview.png", width: 8cm)

      #text(size: 26pt, fill: secondary, weight: "bold")[
        #titolo
      ]

      #v(2cm)
    ]
  ]
  #align(center)[
    #box(
      stroke: 1pt + gray,
      fill: luma(98%),
      inset: 8pt,
      radius: 6pt,
    )[
      #text(weight: "bold", fill: primary)[Stato:] #stato
      #text(weight: "bold", fill: primary)[Versione:] #versione
      #text(weight: "bold", fill: primary)[Presenze:] #presenze

    ]
  ]

  #pagebreak()
  #counter(page).update(1)

  #set page(
    header: [
      #set text(9pt, fill: gray)
      #grid(
        columns: (1fr, 1fr),
        align(left)[GlitchHub Team], align(right)[#titolo],
      )
      #line(length: 100%, stroke: 0.5pt + gray)
    ],
    footer: [
      #set text(12pt)
      #line(length: 100%, stroke: black)
      #v(0.5em)
      #align(center)[
        #context [
          Pag. #counter(page).display() di #counter(page).final().at(0)
        ]
      ]
    ],
  )

  #v(1em)
  #block[
    #text(size: 12pt, weight: "bold")[
      #h(0.5em)
      #heading(numbering: none)[Registro Modifiche]
      #h(0.5em)
    ]
  ]
  #v(1.5em)
  #body
]
