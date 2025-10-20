#import "../../template.typ" : report

#show: report.with(
  titolo: "Verbale esterno 12",
  stato: "Approvato",
  versione: "3",
  presenze: "",
)

  #align(center)[
    #table(
    columns: (auto, auto, auto, auto, auto),
    [*Ver.*], [*Data*], [*Autore*], [*Verificatore*], [*Contenuto*],
    [1.0], [17-10-2025], [Alessandro], [Michele], [Prima riunione]
    )
  ]