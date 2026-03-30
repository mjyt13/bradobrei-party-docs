import { apiRequest } from '../client'
import type {
  LoginRequestDto,
  LoginResponseDto,
  RegisterRequestDto,
  UserResponseDto,
} from '../../types/dto/auth'

export const authService = {
  login(payload: LoginRequestDto) {
    return apiRequest<LoginResponseDto>('/auth/login', {
      method: 'POST',
      auth: false,
      body: payload,
    })
  },
  register(payload: RegisterRequestDto) {
    return apiRequest<{ user: UserResponseDto['user'] }>('/auth/register', {
      method: 'POST',
      auth: false,
      body: payload,
    })
  },
  getMe() {
    return apiRequest<UserResponseDto>('/me')
  },
}
