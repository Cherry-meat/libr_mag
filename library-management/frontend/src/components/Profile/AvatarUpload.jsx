import React, { useState } from 'react'
import api from '../../services/api'
import { useAuth } from '../../contexts/AuthContext'
import { toast } from 'react-toastify'
import { FaUser } from 'react-icons/fa'

function AvatarUpload() {
  const { user, setUser } = useAuth()
  const [uploading, setUploading] = useState(false)

  const handleFileChange = async (e) => {
    const file = e.target.files[0]
    if (!file) return

    const formData = new FormData()
    formData.append('avatar', file)

    setUploading(true)
    try {
      const response = await api.post('/profile/avatar', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      setUser({ ...user, avatar_url: response.data.avatar_url })
      toast.success('Avatar updated successfully!')
    } catch (error) {
      toast.error(error.response?.data?.error || 'Failed to upload avatar')
    } finally {
      setUploading(false)
    }
  }

  return (
    <div className="avatar-upload">
      <div className="avatar-container">
        {user?.avatar_url ? (
          <img src={user.avatar_url} alt="Avatar" className="avatar" />
        ) : (
          <div className="avatar-placeholder">
            <FaUser size={40} />
          </div>
        )}
        <label className="upload-btn">
          <input
            type="file"
            accept="image/*"
            onChange={handleFileChange}
            disabled={uploading}
            hidden
          />
          {uploading ? 'Uploading...' : 'Change Avatar'}
        </label>
      </div>
    </div>
  )
}

export default AvatarUpload