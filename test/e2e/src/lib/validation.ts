
import { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';


export function validateEmptyRoom(room: any) {

    // Validate room parameters
    expect(room.id, 'room id').to.be.a('string')
    expect(room.id.length, 'room id length').to.be.greaterThan(0)

    expect(room.createdAt, 'room created at').to.be.a('string')
    expect(room.createdAt.length, 'room created at length').to.be.greaterThan(0)

    expect(room.state, 'room state').to.be.equal('WAITING_FOR_PLAYERS')

    expect(room.name, 'room name').to.be.a('string')
    expect(room.name.length, 'room name length').to.be.greaterThan(0)

    // Validate empty rounds
    expect(room.currentRound, 'room current round').to.be.null
    expect(room.roundHistory, 'room round history').to.be.a('array')
    expect(room.roundHistory.length, 'room round history length').to.be.equal(0)

    // Validate empty teams
    expect(room.redTeam, 'room red team').to.be.a('array')
    expect(room.redTeam.length, 'room red team length').to.be.equal(0)
    expect(room.blueTeam, 'room blue team').to.be.a('array')
    expect(room.blueTeam.length, 'room blue team length').to.be.equal(0)

}