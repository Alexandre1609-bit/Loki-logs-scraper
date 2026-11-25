# Loki Logs Scraper

Outil en Go qui récupère en continu les logs d'un cluster Kubernetes depuis **Loki** et **pseudonymise les adresses IPv4** avant de les afficher.

## Pourquoi

Au sens du RGPD, une adresse IP est une donnée personnelle. Or les logs applicatifs en sont remplis (`peer.address`, `client_ip`, en-têtes HTTP…). L'objectif est de pouvoir **exploiter les logs** (débogage, analyse, partage avec une équipe) **sans exposer les adresses réelles**.

## Pseudonymisation plutôt qu'anonymisation

L'outil **pseudonymise** les adresses, il ne les anonymise pas. La différence est importante :

|                                   | Anonymisation                     | Pseudonymisation (ce projet)             |
| --------------------------------- | --------------------------------- | ---------------------------------------- |
| Même IP - même valeur             | Non                               | **Oui**                                  |
| Corrélation possible entre lignes | Non                               | **Oui**                                  |
| Retour à l'IP d'origine           | Impossible                        | Possible avec la table de correspondance |
| Statut RGPD                       | N'est plus une donnée personnelle | **Reste une donnée personnelle**         |

Ce choix est volontaire. Un alias stable permet de suivre le parcours d'un client dans les logs (« cette adresse a appelé X puis Y ») sans jamais voir son adresse réelle. En contrepartie, la table de correspondance doit être protégée comme une donnée sensible. Ici, elle ne vit qu'en mémoire et n'est jamais écrite sur disque.

Les alias sont pris dans **`192.0.2.0/24`**, une plage réservée à la documentation ([RFC 5737](https://www.rfc-editor.org/rfc/rfc5737)) : un alias ne peut donc jamais être confondu avec une machine réelle.

## Fonctionnement

1. Toutes les 5 secondes, l'outil interroge l'endpoint `query_range` de Loki sur la fenêtre de temps écoulée depuis l'appel précédent.
2. Chaque ligne de log est analysée et toute adresse IPv4 valide est remplacée par son alias :
   - les zéros en tête sont normalisés (`010.000.000.001` et `10.0.0.1` sont la même adresse, donc le même alias) ;
   - les chaînes qui ressemblent à une IP sans en être une (`999.999.999.999`) sont laissées intactes.

### Exemple

Les adresses ci-dessous sont des exemples. Le format est celui d'un vrai log gRPC.

Avant :

```text
level=info msg="finished call" grpc.method=Check peer.address="10.0.3.17:39224"
level=info msg="finished call" grpc.method=Check peer.address="10.0.5.42:41870"
level=info msg="finished call" grpc.method=Check peer.address="10.0.3.17:39310"
```

Après :

```text
level=info msg="finished call" grpc.method=Check peer.address="192.0.2.1:39224"
level=info msg="finished call" grpc.method=Check peer.address="192.0.2.2:41870"
level=info msg="finished call" grpc.method=Check peer.address="192.0.2.1:39310"
```

`10.0.3.17` apparaît deux fois et reçoit les deux fois l'alias `192.0.2.1`.

## Utilisation

Prérequis : Go 1.26+ et une instance Loki accessible.

```bash
# Loki en local (valeur par défaut : http://localhost:3100/loki/api/v1/query_range)
go run .

# Autre instance : passer l'URL complète de l'endpoint query_range
go run . -loki http://loki.example.local:3100/loki/api/v1/query_range
```

La requête LogQL utilisée récupère tous les namespaces sauf celui de Loki lui-même :

```logql
{namespace=~".+", namespace!="loki"}
```

### Tests

```bash
go test -v ./...
```

Les tests sont écrits en _table-driven_ et couvrent : IP simple, IP au milieu d'un texte, plusieurs IP dans une ligne, IP invalide, texte sans IP, normalisation des zéros en tête.

## Limites connues

- **IPv6 non gérée** : les adresses IPv6 (y compris celles dérivées de l'adresse MAC, format EUI-64) passent en clair.
- **254 adresses au maximum** : au-delà, les alias générés deviennent invalides (`192.0.2.255`, l'adresse de broadcast, puis `192.0.2.256`…).
- **Table en mémoire uniquement** : elle est perdue au redémarrage (les alias changent) et grossit sans limite.
- **Requête LogQL écrite en dur** dans le code.
- **Pas de pagination** : Loki renvoie 100 lignes par requête par défaut, donc des lignes peuvent être perdues en cas de fort volume.

## Pistes d'évolution

- Gestion de l'IPv6 (alias dans `2001:db8::/32`, RFC 3849)
- Gestion du dépassement des 254 adresses
- Encapsulation de l'état dans une struct `Masker` (suppression des variables globales)
- Configuration de la requête et de l'intervalle par fichier ou flags
- `docker-compose` de démonstration (Loki + générateur de logs)
