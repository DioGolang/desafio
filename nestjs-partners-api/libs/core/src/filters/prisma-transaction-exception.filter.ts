import { ExceptionFilter, Catch, ArgumentsHost } from '@nestjs/common';

@Catch()
export class PrismaTransactionExceptionFilter implements ExceptionFilter {
    catch(exception: any, host: ArgumentsHost) {
        const ctx = host.switchToHttp();
        const response = ctx.getResponse();
        const request = ctx.getRequest();

        const status = 400;

        response.status(status).json({
            message: exception.message,
            timestamp: new Date().toISOString(),
            path: request.url,
        });
    }
}
