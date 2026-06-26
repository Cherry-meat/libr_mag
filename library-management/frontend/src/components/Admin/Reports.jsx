import React, { useState } from 'react'
import api from '../../services/api'
import { toast } from 'react-toastify'
import { FaFileDownload } from 'react-icons/fa'

function Reports() {
  const [loading, setLoading] = useState(false)

  const downloadReport = async (type) => {
    setLoading(true)
    try {
      const response = await api.get(`/admin/reports?type=${type}`, {
        responseType: 'blob'
      })
      
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', `${type}_report.csv`)
      document.body.appendChild(link)
      link.click()
      link.remove()
      
      toast.success('Report downloaded successfully')
    } catch (error) {
      toast.error('Failed to download report')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="reports">
      <h2>Reports</h2>
      <div className="reports-grid">
        <div className="report-card" onClick={() => downloadReport('loans')}>
          <FaFileDownload size={40} />
          <h3>Loans Report</h3>
          <p>Download complete loans history</p>
        </div>
        
        <div className="report-card" onClick={() => downloadReport('users')}>
          <FaFileDownload size={40} />
          <h3>Users Report</h3>
          <p>Download users statistics</p>
        </div>
        
        <div className="report-card" onClick={() => downloadReport('books')}>
          <FaFileDownload size={40} />
          <h3>Books Report</h3>
          <p>Download books circulation data</p>
        </div>
      </div>
      {loading && <div className="loading-overlay">Generating report...</div>}
    </div>
  )
}

export default Reports