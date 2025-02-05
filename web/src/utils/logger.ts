import { AxiosRequestConfig, AxiosResponse } from 'axios';

type LogLevel = 'info' | 'warn' | 'error' | 'debug';

interface LogEntry {
  timestamp: string;
  level: LogLevel;
  message: string;
  data?: any;
}

class Logger {
  private static instance: Logger;
  private isDebug: boolean;

  private constructor() {
    this.isDebug = import.meta.env.MODE === 'development';
  }

  public static getInstance(): Logger {
    if (!Logger.instance) {
      Logger.instance = new Logger();
    }
    return Logger.instance;
  }

  private formatLog(level: LogLevel, message: string, data?: any): LogEntry {
    return {
      timestamp: new Date().toISOString(),
      level,
      message,
      data,
    };
  }

  private log(level: LogLevel, message: string, data?: any) {
    const logEntry = this.formatLog(level, message, data);

    switch (level) {
      case 'info':
        console.log(`[${logEntry.timestamp}] [${level.toUpperCase()}] ${message}`, data || '');
        break;
      case 'warn':
        console.warn(`[${logEntry.timestamp}] [${level.toUpperCase()}] ${message}`, data || '');
        break;
      case 'error':
        console.error(`[${logEntry.timestamp}] [${level.toUpperCase()}] ${message}`, data || '');
        break;
      case 'debug':
        if (this.isDebug) {
          console.debug(`[${logEntry.timestamp}] [${level.toUpperCase()}] ${message}`, data || '');
        }
        break;
    }

    return logEntry;
  }

  public info(message: string, data?: any) {
    return this.log('info', message, data);
  }

  public warn(message: string, data?: any) {
    return this.log('warn', message, data);
  }

  public error(message: string, data?: any) {
    return this.log('error', message, data);
  }

  public debug(message: string, data?: any) {
    return this.log('debug', message, data);
  }

  public logRequest(config: AxiosRequestConfig) {
    this.debug('API Request', {
      url: config.url,
      method: config.method,
      headers: config.headers,
      params: config.params,
      data: config.data,
    });
  }

  public logResponse(response: AxiosResponse) {
    this.debug('API Response', {
      url: response.config.url,
      status: response.status,
      statusText: response.statusText,
      headers: response.headers,
      data: response.data,
    });
  }

  public logError(error: any) {
    this.error('API Error', {
      message: error.message,
      code: error.code,
      config: error.config,
      response: error.response?.data,
    });
  }
}

export const logger = Logger.getInstance();