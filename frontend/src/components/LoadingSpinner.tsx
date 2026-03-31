/**
 * LoadingSpinner — Reusable loading indicator component.
 */
function LoadingSpinner({ message = 'Loading...' }: { message?: string }) {
  return (
    <div style={{ textAlign: 'center', padding: '2rem' }}>
      <div
        style={{
          display: 'inline-block',
          width: '2rem',
          height: '2rem',
          border: '3px solid #e0e0e0',
          borderTopColor: '#1a73e8',
          borderRadius: '50%',
          animation: 'spin 0.8s linear infinite',
        }}
      />
      <p style={{ marginTop: '0.5rem', color: '#666' }}>{message}</p>
      <style>{`
        @keyframes spin {
          to { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  )
}

export default LoadingSpinner
