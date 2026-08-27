import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useNavigate, Link, useSearchParams } from 'react-router-dom';
import { registerSchema, RegisterFormData } from '../schemas';
import { authApi } from '../api/auth';
import { invitationsApi } from '../api/invitations';
import { useAuthStore } from '../store/authStore';

export function Register() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { setTokens, setUser, setError: setAuthError } = useAuthStore();
  const [apiError, setApiError] = useState<string>('');
  const [invitationValid, setInvitationValid] = useState<boolean>(false);
  const [invitationError, setInvitationError] = useState<string>('');
  const [checkingInvitation, setCheckingInvitation] = useState(false);

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      invitationCode: searchParams.get('code') || '',
    },
  });

  const invitationCode = watch('invitationCode');

  // Validate invitation code whenever it changes
  const validateInvitation = async (code: string) => {
    if (!code) {
      setInvitationError('');
      setInvitationValid(false);
      return;
    }

    setCheckingInvitation(true);
    try {
      await invitationsApi.validateInvitation(code);
      setInvitationValid(true);
      setInvitationError('');
    } catch (error: any) {
      setInvitationError(error.response?.data?.message || 'Invalid invitation code');
      setInvitationValid(false);
    } finally {
      setCheckingInvitation(false);
    }
  };

  // Debounce validation
  useEffect(() => {
    const timer = setTimeout(() => {
      validateInvitation(invitationCode || '');
    }, 500);

    return () => clearTimeout(timer);
  }, [invitationCode]);

  const onSubmit = async (data: RegisterFormData) => {
    try {
      setApiError('');
      setAuthError(null);

      // Register user
      await authApi.register({
        username: data.username,
        email: data.email,
        password: data.password,
        invitation_code: data.invitationCode,
      } as any);

      // Auto-login after registration
      const authResponse = await authApi.login(data.username, data.password);
      setTokens(authResponse);
      setUser({
        id: 0,
        username: data.username,
        email: data.email,
        role: 'user',
      });

      navigate('/dashboard');
    } catch (error: any) {
      const message = error.response?.data?.message || 'Registration failed. Please try again.';
      setApiError(message);
      setAuthError(message);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 flex items-center justify-center p-4">
      <div className="w-full max-w-md bg-slate-800 rounded-lg shadow-xl p-8">
        <h1 className="text-3xl font-bold text-white mb-2 text-center">Kexel</h1>
        <p className="text-slate-400 text-center mb-8">Create Your Account</p>

        {apiError && (
          <div className="mb-6 p-4 bg-red-900/30 border border-red-700 rounded-lg text-red-200 text-sm">
            {apiError}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Username</label>
            <input
              {...register('username')}
              type="text"
              className="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Choose a username"
            />
            {errors.username && (
              <p className="mt-1 text-sm text-red-400">{errors.username.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Email</label>
            <input
              {...register('email')}
              type="email"
              className="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Enter your email"
            />
            {errors.email && (
              <p className="mt-1 text-sm text-red-400">{errors.email.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Password</label>
            <input
              {...register('password')}
              type="password"
              className="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Create a password"
            />
            {errors.password && (
              <p className="mt-1 text-sm text-red-400">{errors.password.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Confirm Password</label>
            <input
              {...register('confirmPassword')}
              type="password"
              className="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Confirm your password"
            />
            {errors.confirmPassword && (
              <p className="mt-1 text-sm text-red-400">{errors.confirmPassword.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">
              Invitation Code <span className="text-slate-400">(Optional)</span>
            </label>
            <input
              {...register('invitationCode')}
              type="text"
              className="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Paste your invitation code"
            />
            {checkingInvitation && (
              <p className="mt-1 text-sm text-yellow-400">Validating...</p>
            )}
            {invitationCode && invitationValid && (
              <p className="mt-1 text-sm text-green-400">✓ Invitation code is valid</p>
            )}
            {invitationCode && invitationError && (
              <p className="mt-1 text-sm text-red-400">{invitationError}</p>
            )}
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-slate-600 text-white font-semibold py-2 px-4 rounded-lg transition-colors duration-200"
          >
            {isSubmitting ? 'Registering...' : 'Register'}
          </button>
        </form>

        <div className="mt-6 text-center">
          <p className="text-slate-400 text-sm">
            Already have an account?{' '}
            <Link to="/login" className="text-blue-400 hover:text-blue-300 font-medium">
              Login here
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
