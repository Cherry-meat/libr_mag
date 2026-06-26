import React from 'react'
import { useAuth } from '../../contexts/AuthContext'
import { FaTicketAlt, FaQrcode, FaUser, FaEnvelope } from 'react-icons/fa'

function TicketView() {
  const { user } = useAuth()

  if (!user?.ticket_linked) {
    return (
      <div className="ticket-not-linked">
        <FaTicketAlt size={50} />
        <h2>No Ticket Linked</h2>
        <p>Please link your plastic ticket in your profile</p>
        <button onClick={() => window.location.href = '/profile'}>
          Go to Profile
        </button>
      </div>
    )
  }

  return (
    <div className="ticket-container">
      <div className="ticket-card">
        <div className="ticket-header">
          <h2>Electronic Reader Ticket</h2>
          <div className="ticket-badge">
            {user.is_blocked ? 'BLOCKED' : 'ACTIVE'}
          </div>
        </div>
        <div className="ticket-body">
          <div className="ticket-qr">
            <FaQrcode size={100} />
            <p>{user.ticket_number}</p>
          </div>
          <div className="ticket-info">
            <div className="info-row">
              <FaUser />
              <span>{user.first_name} {user.last_name}</span>
            </div>
            <div className="info-row">
              <FaEnvelope />
              <span>{user.email}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default TicketView