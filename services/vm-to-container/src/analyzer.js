/**
 * VMAnalyzer - Analyzes VMs for containerization
 */

const fs = require('fs-extra');
const path = require('path');
const { execSync } = require('child_process');
const glob = require('glob');

class VMAnalyzer {
  constructor(options = {}) {
    this.options = {
      tempDir: options.tempDir || '/tmp/vm2container',
      ...options,
    };
  }

  /**
   * Analyze a VM source
   */
  async analyze(source, type = 'auto') {
    const detectedType = type === 'auto' ? this.detectType(source) : type;
    
    const analysis = {
      source,
      type: detectedType,
      timestamp: new Date().toISOString(),
      os: {},
      applications: [],
      services: [],
      ports: [],
      files: [],
      users: [],
      packages: [],
      estimatedSize: 0,
      score: 0,
      warnings: [],
      recommendations: [],
    };

    switch (detectedType) {
      case 'vmware':
        await this.analyzeVMware(source, analysis);
        break;
      case 'virtualbox':
        await this.analyzeVirtualBox(source, analysis);
        break;
      case 'raw':
        await this.analyzeRawDisk(source, analysis);
        break;
      case 'ssh':
        await this.analyzeSSH(source, analysis);
        break;
      default:
        throw new Error(`Unsupported VM type: ${detectedType}`);
    }

    // Calculate containerization score
    analysis.score = this.calculateScore(analysis);

    return analysis;
  }

  /**
   * Detect VM type from source
   */
  detectType(source) {
    if (source.endsWith('.vmx')) return 'vmware';
    if (source.endsWith('.vbox')) return 'virtualbox';
    if (source.endsWith('.vmdk') || source.endsWith('.vdi') || source.endsWith('.qcow2')) {
      return 'raw';
    }
    if (source.includes('@') && source.includes(':')) return 'ssh';
    return 'raw';
  }

  /**
   * Analyze VMware VM
   */
  async analyzeVMware(source, analysis) {
    // Parse VMX file
    const vmxPath = source.endsWith('.vmx') ? source : path.join(source, '*.vmx');
    const vmxFiles = glob.sync(vmxPath);

    if (vmxFiles.length === 0) {
      throw new Error('No VMX file found');
    }

    const vmxContent = await fs.readFile(vmxFiles[0], 'utf8');
    const vmxConfig = this.parseVMX(vmxContent);

    analysis.os.name = vmxConfig.guestOS || 'unknown';
    analysis.os.arch = vmxConfig.guestOS.includes('64') ? 'x86_64' : 'x86';

    // Detect disks
    const disks = Object.keys(vmxConfig)
      .filter(key => key.startsWith('scsi') || key.startsWith('sata') || key.startsWith('ide'))
      .filter(key => key.endsWith('.fileName'))
      .map(key => vmxConfig[key]);

    analysis.files = disks;

    // Estimate size
    for (const disk of disks) {
      const diskPath = path.join(path.dirname(vmxFiles[0]), disk);
      try {
        const stats = await fs.stat(diskPath);
        analysis.estimatedSize += Math.round(stats.size / 1024 / 1024);
      } catch (e) {
        // Disk file might not exist locally
      }
    }

    // Detect common applications
    analysis.applications = this.detectApplications(analysis.os.name);
    analysis.services = this.inferServices(analysis.applications);
    analysis.ports = this.inferPorts(analysis.applications);

    return analysis;
  }

  /**
   * Analyze VirtualBox VM
   */
  async analyzeVirtualBox(source, analysis) {
    const vboxPath = source.endsWith('.vbox') ? source : path.join(source, '*.vbox');
    const vboxFiles = glob.sync(vboxPath);

    if (vboxFiles.length === 0) {
      throw new Error('No VBOX file found');
    }

    const vboxContent = await fs.readFile(vboxFiles[0], 'utf8');
    
    // Simple parsing - in production use proper XML parser
    const osTypeMatch = vboxContent.match(/OSType="([^"]+)"/);
    if (osTypeMatch) {
      analysis.os.name = osTypeMatch[1];
      analysis.os.arch = analysis.os.name.includes('_64') ? 'x86_64' : 'x86';
    }

    // Detect disks
    const diskMatches = vboxContent.matchAll(/location="([^"]+\.(vdi|vmdk|vhd))"/g);
    for (const match of diskMatches) {
      analysis.files.push(match[1]);
    }

    analysis.applications = this.detectApplications(analysis.os.name);
    analysis.services = this.inferServices(analysis.applications);
    analysis.ports = this.inferPorts(analysis.applications);

    return analysis;
  }

  /**
   * Analyze raw disk image
   */
  async analyzeRawDisk(source, analysis) {
    analysis.files = [source];
    
    try {
      const stats = await fs.stat(source);
      analysis.estimatedSize = Math.round(stats.size / 1024 / 1024);
    } catch (e) {
      // Ignore
    }

    // Try to detect OS from disk
    analysis.os.name = 'linux'; // Default assumption
    analysis.os.arch = 'x86_64';

    analysis.applications = this.detectApplications('linux');
    analysis.services = this.inferServices(analysis.applications);
    analysis.ports = this.inferPorts(analysis.applications);

    return analysis;
  }

  /**
   * Analyze VM via SSH
   */
  async analyzeSSH(source, analysis) {
    // SSH analysis would use ssh2 library
    // For now, return basic structure
    analysis.os.name = 'linux';
    analysis.os.arch = 'x86_64';
    analysis.applications = [];
    analysis.services = [];
    analysis.ports = [];

    analysis.warnings.push('SSH-based analysis requires manual service detection');

    return analysis;
  }

  /**
   * Parse VMX file content
   */
  parseVMX(content) {
    const config = {};
    const lines = content.split('\n');

    for (const line of lines) {
      const match = line.match(/^(.+)\s*=\s*"(.+)"$/);
      if (match) {
        config[match[1].trim()] = match[2].trim();
      }
    }

    return config;
  }

  /**
   * Detect applications based on OS
   */
  detectApplications(osName) {
    const apps = [];

    const commonApps = {
      'linux': [
        { name: 'nginx', containerizable: true },
        { name: 'apache', containerizable: true },
        { name: 'mysql', containerizable: true },
        { name: 'postgresql', containerizable: true },
        { name: 'redis', containerizable: true },
        { name: 'nodejs', containerizable: true },
        { name: 'python', containerizable: true },
        { name: 'java', containerizable: true },
      ],
      'windows': [
        { name: 'iis', containerizable: true },
        { name: 'sqlserver', containerizable: true },
        { name: '.net', containerizable: true },
      ],
    };

    const osKey = osName.toLowerCase().includes('win') ? 'windows' : 'linux';
    return commonApps[osKey] || commonApps['linux'];
  }

  /**
   * Infer services from applications
   */
  inferServices(applications) {
    const serviceMap = {
      'nginx': [{ name: 'nginx', ports: [80, 443] }],
      'apache': [{ name: 'apache2', ports: [80, 443] }],
      'mysql': [{ name: 'mysql', ports: [3306] }],
      'postgresql': [{ name: 'postgresql', ports: [5432] }],
      'redis': [{ name: 'redis', ports: [6379] }],
      'nodejs': [{ name: 'node', ports: [3000, 8080] }],
    };

    const services = [];
    for (const app of applications) {
      if (serviceMap[app.name]) {
        services.push(...serviceMap[app.name]);
      }
    }

    return services;
  }

  /**
   * Infer ports from applications
   */
  inferPorts(applications) {
    const ports = [];
    const portMap = {
      'nginx': [80, 443],
      'apache': [80, 443],
      'mysql': [3306],
      'postgresql': [5432],
      'redis': [6379],
      'mongodb': [27017],
      'elasticsearch': [9200],
      'nodejs': [3000, 8080],
    };

    for (const app of applications) {
      if (portMap[app.name]) {
        ports.push(...portMap[app.name]);
      }
    }

    return [...new Set(ports)];
  }

  /**
   * Detect services in VM
   */
  async detectServices(source, type) {
    const analysis = await this.analyze(source, type);
    return analysis.services.map(s => ({
      name: s.name,
      type: 'systemd',
      status: 'unknown',
      ports: s.ports || [],
    }));
  }

  /**
   * Calculate containerization score
   */
  calculateScore(analysis) {
    let score = 50; // Base score

    // Bonus for known OS
    if (analysis.os.name && analysis.os.name !== 'unknown') {
      score += 10;
    }

    // Bonus for containerizable applications
    const containerizableApps = analysis.applications.filter(a => a.containerizable).length;
    score += containerizableApps * 5;

    // Penalty for large disks
    if (analysis.estimatedSize > 10000) { // > 10GB
      score -= 10;
    }

    // Penalty for warnings
    score -= analysis.warnings.length * 5;

    return Math.max(0, Math.min(100, score));
  }
}

module.exports = { VMAnalyzer };
