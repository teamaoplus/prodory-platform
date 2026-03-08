#!/usr/bin/env node

/**
 * VM-to-Container Migration Tool
 * 
 * A CLI tool for migrating virtual machines to Docker containers.
 * Analyzes VM configuration, extracts applications, and generates
 * container images with appropriate Dockerfiles.
 * 
 * Features:
 * - VM analysis and discovery
 * - Application dependency detection
 * - Dockerfile generation
 * - Container image building
 * - Docker Compose generation
 * - Kubernetes manifest generation
 */

const { Command } = require('commander');
const chalk = require('chalk');
const ora = require('ora');
const inquirer = require('inquirer');
const fs = require('fs-extra');
const path = require('path');

const { VMAnalyzer } = require('./analyzer');
const { ContainerBuilder } = require('./builder');
const { ImageGenerator } = require('./generator');
const { MigrationEngine } = require('./engine');

const program = new Command();
const version = '1.0.0';

program
  .name('vm2container')
  .description('Virtual Machine to Container Migration Tool')
  .version(version)
  .configureOutput({
    outputError: (str, write) => write(chalk.red(str))
  });

program
  .command('analyze')
  .description('Analyze a VM for containerization')
  .option('-s, --source <path>', 'VM source path or connection string')
  .option('-t, --type <type>', 'VM type (vmware, virtualbox, raw, ssh)', 'auto')
  .option('-o, --output <file>', 'Output file for analysis report', 'analysis.json')
  .action(async (options) => {
    console.log(chalk.cyan(`
╔╗ ╦  ╦╔╦╗╔═╗╔═╗╔═╗  ╔╦╗╔═╗╔╦╗╔═╗╔═╗╦═╗╔═╗
╠╩╗║  ║║║║║╣ ║╣ ╚═╗   ║║╠═╣ ║ ║╣ ║ ╦╠╦╝╚═╗
╚═╝╩═╝╩╩ ╩╚═╝╚═╝╚═╝  ═╩╝╩ ╩ ╩ ╚═╝╚═╝╩╚═╚═╝
    `));

    try {
      const source = options.source || await promptForSource();
      const spinner = ora('Analyzing VM...').start();

      const analyzer = new VMAnalyzer();
      const report = await analyzer.analyze(source, options.type);

      spinner.succeed('Analysis complete');

      // Save report
      await fs.writeJson(options.output, report, { spaces: 2 });
      console.log(chalk.green(`✓ Analysis report saved to ${options.output}`));

      // Display summary
      displayAnalysisSummary(report);

    } catch (error) {
      console.error(chalk.red(`Error: ${error.message}`));
      process.exit(1);
    }
  });

program
  .command('migrate')
  .description('Migrate VM to container')
  .option('-s, --source <path>', 'VM source path')
  .option('-t, --type <type>', 'VM type', 'auto')
  .option('-n, --name <name>', 'Container image name')
  .option('-o, --output <dir>', 'Output directory', './output')
  .option('--base-image <image>', 'Base Docker image')
  .option('--expose <ports>', 'Ports to expose (comma-separated)')
  .option('--env <envs>', 'Environment variables (KEY=value,...)')
  .option('--volumes <volumes>', 'Volumes to mount')
  .option('--generate-compose', 'Generate Docker Compose file', false)
  .option('--generate-k8s', 'Generate Kubernetes manifests', false)
  .option('--build', 'Build container image', false)
  .option('--push', 'Push image to registry', false)
  .action(async (options) => {
    console.log(chalk.cyan('\n🚀 Starting VM to Container Migration\n'));

    try {
      const source = options.source || await promptForSource();
      const imageName = options.name || await promptForImageName();

      const spinner = ora('Initializing migration engine...').start();
      const engine = new MigrationEngine({
        outputDir: options.output,
        baseImage: options.baseImage,
        ports: options.expose ? options.expose.split(',').map(p => p.trim()) : [],
        env: parseEnv(options.env),
        volumes: options.volumes ? options.volumes.split(',').map(v => v.trim()) : [],
      });
      spinner.succeed('Migration engine ready');

      // Step 1: Analyze
      spinner.start('Analyzing VM...');
      const analysis = await engine.analyze(source, options.type);
      spinner.succeed('VM analysis complete');

      // Step 2: Generate Dockerfile
      spinner.start('Generating Dockerfile...');
      const dockerfile = await engine.generateDockerfile(analysis, imageName);
      spinner.succeed('Dockerfile generated');

      // Step 3: Generate supporting files
      spinner.start('Generating container configuration...');
      await engine.generateConfigFiles(analysis, imageName, {
        compose: options.generateCompose,
        kubernetes: options.generateK8s,
      });
      spinner.succeed('Configuration files generated');

      // Step 4: Build image (optional)
      if (options.build) {
        spinner.start('Building container image...');
        const image = await engine.buildImage(imageName, options.output);
        spinner.succeed(`Image built: ${image}`);

        if (options.push) {
          spinner.start('Pushing image to registry...');
          await engine.pushImage(image);
          spinner.succeed('Image pushed');
        }
      }

      // Summary
      console.log(chalk.green('\n✓ Migration completed successfully!\n'));
      console.log('Output directory:', chalk.cyan(options.output));
      console.log('Dockerfile:', chalk.cyan(path.join(options.output, 'Dockerfile')));
      if (options.generateCompose) {
        console.log('Docker Compose:', chalk.cyan(path.join(options.output, 'docker-compose.yml')));
      }
      if (options.generateK8s) {
        console.log('K8s manifests:', chalk.cyan(path.join(options.output, 'k8s/')));
      }
      console.log();

    } catch (error) {
      console.error(chalk.red(`\n✗ Migration failed: ${error.message}`));
      process.exit(1);
    }
  });

program
  .command('generate')
  .description('Generate Dockerfile from analysis')
  .option('-a, --analysis <file>', 'Analysis report file', 'analysis.json')
  .option('-n, --name <name>', 'Image name', 'migrated-app')
  .option('-o, --output <dir>', 'Output directory', './output')
  .option('--base-image <image>', 'Base Docker image')
  .action(async (options) => {
    try {
      const spinner = ora('Loading analysis...').start();
      const analysis = await fs.readJson(options.analysis);
      spinner.succeed('Analysis loaded');

      const generator = new ImageGenerator({
        outputDir: options.output,
        baseImage: options.baseImage,
      });

      spinner.start('Generating Dockerfile...');
      await generator.generate(analysis, options.name);
      spinner.succeed('Dockerfile generated');

      console.log(chalk.green(`\n✓ Dockerfile saved to ${options.output}/Dockerfile\n`));

    } catch (error) {
      console.error(chalk.red(`Error: ${error.message}`));
      process.exit(1);
    }
  });

program
  .command('build')
  .description('Build container image from generated files')
  .option('-d, --directory <dir>', 'Build context directory', './output')
  .option('-n, --name <name>', 'Image name', 'migrated-app')
  .option('-t, --tag <tag>', 'Image tag', 'latest')
  .option('--push', 'Push after build', false)
  .action(async (options) => {
    try {
      const spinner = ora('Building container image...').start();

      const builder = new ContainerBuilder();
      const imageName = `${options.name}:${options.tag}`;
      
      const image = await builder.build({
        context: options.directory,
        dockerfile: path.join(options.directory, 'Dockerfile'),
        tag: imageName,
      });

      spinner.succeed(`Image built: ${image}`);

      if (options.push) {
        spinner.start('Pushing image...');
        await builder.push(image);
        spinner.succeed('Image pushed');
      }

      console.log(chalk.green('\n✓ Build complete!\n'));

    } catch (error) {
      console.error(chalk.red(`Error: ${error.message}`));
      process.exit(1);
    }
  });

program
  .command('list-services')
  .description('List services detected in VM')
  .option('-s, --source <path>', 'VM source path')
  .option('-t, --type <type>', 'VM type', 'auto')
  .action(async (options) => {
    try {
      const source = options.source || await promptForSource();
      const spinner = ora('Detecting services...').start();

      const analyzer = new VMAnalyzer();
      const services = await analyzer.detectServices(source, options.type);

      spinner.succeed(`Found ${services.length} services`);

      console.log('\nDetected Services:');
      services.forEach(service => {
        console.log(`  • ${chalk.cyan(service.name)} (${service.type})`);
        console.log(`    Status: ${service.status}`);
        if (service.ports.length > 0) {
          console.log(`    Ports: ${service.ports.join(', ')}`);
        }
        console.log();
      });

    } catch (error) {
      console.error(chalk.red(`Error: ${error.message}`));
      process.exit(1);
    }
  });

// Helper functions
async function promptForSource() {
  const { source } = await inquirer.prompt([{
    type: 'input',
    name: 'source',
    message: 'Enter VM source path or connection string:',
    validate: (input) => input.length > 0 || 'Source is required',
  }]);
  return source;
}

async function promptForImageName() {
  const { name } = await inquirer.prompt([{
    type: 'input',
    name: 'name',
    message: 'Enter container image name:',
    default: 'migrated-app',
    validate: (input) => /^[a-z0-9._-]+$/.test(input) || 'Invalid image name',
  }]);
  return name;
}

function parseEnv(envString) {
  if (!envString) return {};
  const env = {};
  envString.split(',').forEach(pair => {
    const [key, value] = pair.split('=');
    if (key && value) {
      env[key.trim()] = value.trim();
    }
  });
  return env;
}

function displayAnalysisSummary(report) {
  console.log(chalk.cyan('\n=== Analysis Summary ===\n'));

  console.log('Operating System:', chalk.yellow(report.os.name, report.os.version));
  console.log('Architecture:', chalk.yellow(report.os.arch));
  console.log();

  console.log('Detected Applications:');
  report.applications.forEach(app => {
    const status = app.containerizable 
      ? chalk.green('✓') 
      : chalk.yellow('⚠');
    console.log(`  ${status} ${app.name} ${app.version || ''}`);
  });
  console.log();

  console.log('Services:', report.services.length);
  console.log('Open Ports:', report.ports.join(', ') || 'None detected');
  console.log('Estimated Image Size:', chalk.yellow(`${report.estimatedSize}MB`));
  console.log('Containerization Score:', chalk.yellow(`${report.score}/100`));
  console.log();

  if (report.warnings.length > 0) {
    console.log(chalk.yellow('Warnings:'));
    report.warnings.forEach(w => console.log(`  ⚠ ${w}`));
    console.log();
  }

  if (report.recommendations.length > 0) {
    console.log('Recommendations:');
    report.recommendations.forEach(r => console.log(`  • ${r}`));
    console.log();
  }
}

program.parse();
