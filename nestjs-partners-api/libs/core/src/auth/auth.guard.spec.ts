import { AuthGuard } from './auth.guard';

describe('AuthGuard', () => {
  it('should be defined', () => {
    // @ts-ignore
    expect(new AuthGuard()).toBeDefined();
  });
});
