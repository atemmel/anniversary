package main

import(
	"time"
	"fmt"
	"math/rand"
	"github.com/bwmarrin/discordgo"
)

type Team int

const(
	TeamBlu = 0
	TeamRed = 1
)

type WhisperParticipant struct {
	id string
	team Team
	turn int
}

func createTeams(ids []string) []WhisperParticipant {
	wp := make([]WhisperParticipant, 0)
	for i := 0; len(ids) != 0 ; i++ {
		index := rand.Intn(len(ids))
		w := WhisperParticipant{
			ids[index],
			-1,
			i,
		}

		if i % 2 == 0 {
			w.team = TeamBlu
		} else {
			w.team = TeamRed
		}

		wp = append(wp, w)
		ids = append(ids[:index], ids[index +1:]...)
	}

	return wp
}

func sendStartMessage(s *discordgo.Session, r1, b1 WhisperParticipant) {
	strs := [...]string{
		"Jag är small brain",
	}

	bchan, _ := s.UserChannelCreate(b1.id)
	rchan, _ := s.UserChannelCreate(r1.id)

	str := "Du börjar viskleken med: \"" + strs[rand.Intn(len(strs))] + "\""

	s.ChannelMessageSend(bchan.ID, str)
	s.ChannelMessageSend(rchan.ID, str)
}

func startWhisperGame(s *discordgo.Session) {
	g, err := s.State.Guild(guildId)
	if err != nil {
		// Could not find guild.
		return
	}

	conns := make([]string, 0)
	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.ChannelID == generalChannelId {
			conns = append(conns, vs.UserID)
		}
	}

	teams := createTeams(conns)

	for i := range teams {
		if teams[i].team == TeamRed {
			s.GuildMemberRoleAdd(guildId, teams[i].id, teamRedId)
		} else if teams[i].team == TeamBlu {
			s.GuildMemberRoleAdd(guildId, teams[i].id, teamBluId)
		}
	}

	for i := range conns {
		assertMove(s, guildId, conns[i], &muteChannelId)
	}

	reds := make([]WhisperParticipant, 0)
	blus := make([]WhisperParticipant, 0)

	for i := range teams {
		if teams[i].team == TeamRed {
			reds = append(reds, teams[i])
		} else if teams[i].team == TeamBlu {
			blus = append(blus, teams[i])
		}
	}

	max := 0
	if len(reds) < len(blus) {
		max = len(blus)
	} else {
		max = len(reds)
	}

	sendStartMessage(s, reds[0], blus[0])

	for i := 0; i < max; i++ {
		var r WhisperParticipant
		var b WhisperParticipant

		if i < len(reds) {
			r = reds[i]
		}
		if i < len(reds) {
			b = blus[i]
		}

		assertMove(s, guildId, b.id, &whisperChannelBluId)
		assertMove(s, guildId, r.id, &whisperChannelRedId)
		time.Sleep(time.Duration(1000 * 1000 * 1000 * 3))
		if i + 1 < len(blus) {
			assertMove(s, guildId, blus[i + 1].id, &whisperChannelBluId)
		}
		if i + 1 < len(reds) {
			assertMove(s, guildId, reds[i + 1].id, &whisperChannelRedId)
		}
		time.Sleep(time.Duration(1000 * 1000 * 1000 * 3))
		assertMove(s, guildId, b.id, &whisperChannelBluId)
		assertMove(s, guildId, r.id, &whisperChannelRedId)
		time.Sleep(time.Duration(1000 * 1000 * 1000 * 3))
	}

	for i := range teams {
		assertMove(s, guildId, teams[i].id, &muteChannelId)
	}

	time.Sleep(time.Duration(1000 * 1000 * 1000 * 3))

	cleanupWhisperGame(s, teams)
}

func assertMove(s *discordgo.Session, guildId, memId string, chanId *string) {
	err := s.GuildMemberMove(guildId, memId, chanId)
	for err != nil {
		err = s.GuildMemberMove(guildId, memId, chanId)
		fmt.Println(err)
	}
}


func cleanupWhisperGame(s *discordgo.Session, teams []WhisperParticipant) {
	for i := range teams {
		if teams[i].team == TeamRed {
			s.GuildMemberRoleRemove(guildId, teams[i].id, teamRedId)
		} else if teams[i].team == TeamBlu {
			s.GuildMemberRoleRemove(guildId, teams[i].id, teamBluId)
		}
		s.GuildMemberMove(guildId, teams[i].id, &generalChannelId)
	}
}

