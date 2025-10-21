#import "../../template.typ": report

#show: report.with(
  titolo: "Verbale esterno 12",
  stato: "Approvato",
  versione: "3",
  presenze: "",
  htmlId: "Candidatura-DocumentiEsterni-VerbaliEsterni",
)

#align(center)[
  #table(
    columns: (auto, auto, auto, auto, auto),
    [*Ver.*], [*Data*], [*Autore*], [*Verificatore*], [*Contenuto*],
    [1.0], [14-10-2025], [Riccardo], [Michele], [Prima riunione],
  )
]
