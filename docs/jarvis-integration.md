# Intégration JARVIS — apple-music-mcp

Guide d'intégration pour agents vocaux et assistants (JARVIS, Cursor, Claude Desktop).

## Vue d'ensemble

`apple-music-mcp` expose Apple Music via le **Model Context Protocol (MCP)**. Toutes les commandes passent par Music.app en local — aucun cloud, aucun token Apple ID.

```
JARVIS / Agent IA
       │
       ▼  MCP (stdio JSON-RPC)
apple-music-mcp
       │
       ▼  Apple Events (osascript)
   Music.app
```

## Installation rapide

```bash
# Compiler depuis les sources
cd apple-music-mcp
go build -o apple-music-mcp ./cmd/apple-music-mcp/
cp apple-music-mcp ~/.local/bin/

# Vérifier
apple-music-mcp doctor
```

### Configuration Cursor / JARVIS

Ajouter dans `~/.cursor/mcp.json` :

```json
{
  "mcpServers": {
    "apple-music": {
      "command": "/Users/TOI/.local/bin/apple-music-mcp",
      "args": ["serve"]
    }
  }
}
```

Recharger les serveurs MCP après modification.

## Permissions macOS

Au premier appel, macOS demande l'autorisation **Automation** pour contrôler Music.app :

**Réglages Système → Confidentialité et sécurité → Automatisation**

## Variables d'environnement

| Variable | Défaut | Description |
|----------|--------|-------------|
| `APPLE_MUSIC_MCP_READ_ONLY` | `0` | `1` = lecture seule (pas de play/delete/favorite) |
| `APPLE_MUSIC_MCP_VERBOSE` | `0` | `1` = logs debug sur stderr |
| `APPLE_MUSIC_MCP_LIVE_TESTS` | `0` | `1` = active `test-live` |
| `APPLE_MUSIC_MCP_LIVE_SEARCH` | `Werenoi` | Artiste/requête pour les tests live |
| `APPLE_MUSIC_MCP_LIVE_PLAYLIST` | `Jarvis MCP Live Test` | Nom de playlist temporaire pour les tests |
| `APPLE_MUSIC_MCP_LIVE_KEEP_PLAYLIST` | `0` | `1` = ne pas supprimer la playlist de test |

## Outils MCP disponibles

### `music_get_state`

État de lecture actuel (piste, volume, shuffle, repeat).

```json
{}
```

**Cas d'usage JARVIS :** *« Qu'est-ce qui joue ? »*, *« Quelle chanson c'est ? »*

---

### `music_search`

Recherche dans la **bibliothèque locale** (pas le catalogue Apple Music streaming).

```json
{
  "query": "Werenoi",
  "types": ["track"],
  "limit": 10
}
```

**Retour :** `tracks[]` avec `persistent_id`, `name`, `artist`, `album`, `favorited`.

**Cas d'usage :** *« Cherche Werenoi dans ma bibliothèque »*

---

### `music_playback`

Contrôle la lecture.

| Action | Paramètres | Description |
|--------|------------|-------------|
| `play_track` | `target_id`, `target_type: "track"` | Lire un morceau par ID persistant |
| `play_album` | `target_id`, `target_type: "album"` | Lire un album |
| `play_artist` | `target_id`, `target_type: "artist"` | Lire un artiste |
| `play_playlist` | `target_id`, `target_type: "playlist"` | Lire une playlist |
| `play_url` | `target_id` (URL Apple Music) | Ouvrir une URL Apple Music |
| `play` / `pause` / `toggle` / `stop` | — | Contrôles basiques |
| `next` / `previous` | — | Piste suivante/précédente |
| `seek_absolute` | `seek_position` (secondes) | Aller à une position |

**Exemple — lancer Werenoi :**

```json
{
  "action": "play_track",
  "target_id": "719E704A4F635085",
  "target_type": "track"
}
```

**Workflow JARVIS recommandé :**

1. `music_search` → récupérer `persistent_id`
2. `music_playback` avec `play_track`
3. `music_get_state` → confirmer

---

### `music_favorites`

Gestion des favoris et notes.

| Action | Paramètres |
|--------|------------|
| `favorite` | `target_id`, `target_type: "track"` |
| `unfavorite` | `target_id`, `target_type` |
| `set_rating` | `target_id`, `rating` (0-100, 100 = 5★) |
| `get` | `target_id` |
| `list` | — (liste les favoris) |

**Exemple :**

```json
{
  "action": "favorite",
  "target_id": "719E704A4F635085",
  "target_type": "track"
}
```

---

### `music_playlists`

Gestion des playlists.

| Action | Paramètres |
|--------|------------|
| `list` | — |
| `get` | `playlist_id` |
| `create` | `new_name` |
| `delete` | `playlist_id` |
| `add_tracks` | `playlist_id`, `track_ids[]` |
| `remove_tracks` | `playlist_id`, `track_ids[]` |
| `rename` | `playlist_id`, `new_name` |

**Exemple — créer la playlist Jarvis :**

```json
{ "action": "create", "new_name": "Jarvis" }
```

Puis ajouter un morceau :

```json
{
  "action": "add_tracks",
  "playlist_id": "A789243C358C7798",
  "track_ids": ["719E704A4F635085"]
}
```

---

### `music_preferences`

Volume, shuffle, repeat, sorties AirPlay.

```json
{ "action": "set_volume", "volume": 40 }
{ "action": "set_shuffle", "shuffle_mode": "songs" }
{ "action": "set_repeat", "repeat_mode": "all" }
```

---

### `music_doctor` / `music_capabilities`

Diagnostics et matrice de capacités. À appeler en cas d'erreur ou au démarrage de JARVIS.

## Scénarios JARVIS (prompts → outils)

| Demande utilisateur | Séquence d'outils |
|---------------------|-------------------|
| *« Mets du Werenoi »* | `music_search` → `music_playback(play_track)` |
| *« Mets en favori »* | `music_get_state` → `music_favorites(favorite)` |
| *« Crée une playlist Jarvis »* | `music_playlists(create)` → `music_playlists(add_tracks)` |
| *« Pause »* | `music_playback(pause)` |
| *« Volume à 30 % »* | `music_preferences(set_volume)` |
| *« Quelles sont mes playlists ? »* | `music_playlists(list)` |

## Tests live (obligatoire avant déploiement JARVIS)

La suite live simule le workflow complet agent :

```bash
APPLE_MUSIC_MCP_LIVE_TESTS=1 apple-music-mcp test-live
```

**Étapes exécutées :**

1. Vérification automation Music.app
2. Recherche bibliothèque (défaut : Werenoi)
3. Lecture via `play_track`
4. Lecture de l'état
5. Ajout aux favoris
6. Création playlist temporaire + ajout du morceau
7. Nettoyage de la playlist de test

**Sortie attendue :** toutes les lignes en ✅, `0 failed`.

Pour conserver la playlist de test :

```bash
APPLE_MUSIC_MCP_LIVE_TESTS=1 \
APPLE_MUSIC_MCP_LIVE_PLAYLIST="Jarvis" \
# (modifier livetest pour KeepPlaylist - not exposed yet)
```

## Limitations connues

| Fonctionnalité | Statut | Alternative |
|----------------|--------|-------------|
| File d'attente (Up Next) | ❌ Non supporté | — |
| Recherche catalogue Apple Music | ❌ Bibliothèque locale uniquement | `play_url` avec URL Apple Music |
| Autoplay | ❌ Non scriptable | — |
| Pistes streaming sans métadonnées | ⚠️ Parfois | Utiliser morceaux de la bibliothèque |

## Dépannage

| Symptôme | Cause probable | Solution |
|----------|----------------|----------|
| `success: false` sur tous les outils JSON | Permission automation | `doctor` + Réglages Système |
| `music_search` retourne 0 résultats | Artiste absent de la bibliothèque | `play_url` ou ajouter à la bibliothèque |
| `music_get_state` sans piste | Lecture streaming catalogue | `play_track` avec ID bibliothèque |
| Serveur MCP non visible | Config non rechargée | Redémarrer Cursor / recharger MCP |

## Sécurité

- Aucun secret stocké
- Aucune télémétrie
- stdout réservé au JSON-RPC MCP (logs sur stderr uniquement)
- Mode lecture seule disponible : `APPLE_MUSIC_MCP_READ_ONLY=1`

## Références

- [README](../README.md) — installation et architecture
- [capability-matrix.md](capability-matrix.md) — matrice complète des capacités
- [Model Context Protocol](https://modelcontextprotocol.io)
