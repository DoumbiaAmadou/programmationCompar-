Le programme est effectu� en java pour des raisons de connaissances du langage.

Il fait environ 90 lignes (en comptant les espaces)

Chaque mot dans le dictionnaire est stock� dans un sous dictionnaire correspondant � sa taille sous la forme d'un Element, qui contient le mot en brute, et sous sa forme tri�e alphab�tiquement.
Lorsqu'un mot est recherch�, le sous dictionnaire de taille correspondante est parcouru, avec une recherche sur le mot tri�. Ceci permet d'am�liorer grandement le temps de comparaison.


La complexit� de la recherche est de l'ordre O(n).

