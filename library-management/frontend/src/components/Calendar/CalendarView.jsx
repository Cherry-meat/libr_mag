import React, { useState, useEffect } from 'react'
import api from '../../services/api'
import { toast } from 'react-toastify'

function CalendarView() {
  const [events, setEvents] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchCalendarData()
  }, [])

  const fetchCalendarData = async () => {
    try {
      const response = await api.get('/calendar')
      setEvents(response.data)
    } catch (error) {
      toast.error('Failed to fetch calendar data')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="loading">Loading calendar...</div>
  }

  return (
    <div className="calendar-view">
      <h2>Book Loan Calendar</h2>
      {events.length === 0 ? (
        <p>No events found</p>
      ) : (
        <div className="events-list">
          {events.map((event, index) => (
            <div key={index} className="event-item">
              <h4>{event.title}</h4>
              <p>Borrowed: {new Date(event.loan_date).toLocaleDateString()}</p>
              <p>Due: {new Date(event.due_date).toLocaleDateString()}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default CalendarView