import { api } from './index'

export const chatApi = {
  send: (message: string, sessionId: string) =>
    api.post('/chat', { message, session_id: sessionId }).then(r => r.data),
}