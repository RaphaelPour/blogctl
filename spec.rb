require 'rspec'
require 'open3'
require 'tmpdir'

def blogctl(args)
  go_args = "ci-build/blogctl.test -test.run '^Test_'"\
    " -test.coverprofile=coverage/#{rand 0..10_000}.out"

  Open3.capture3({'TEST_ARGS' =>  args}, go_args)
end

describe 'CLI' do

  context 'Trivial test cases' do
    it 'shows a version' do
      out, err, _ = blogctl("--version")
      expect(err).to eq ""
      expect(out).to match(/BuildVersion:/)
      expect(out).to match(/BuildDate:/)
    end

    it 'shows general help without any args' do
      out,err,_ = blogctl("")
      expect(err).to eq ""
      expect(out).to include 'Blogctl manages blog markdown-based posts'\
                             ' database-less and generates a static website on-demand'
      expect(out).to include 'Usage'
      expect(out).to include 'Available Commands'
      expect(out).to include 'Use "blogctl [command] --help" for more information about a command.'
    end

    it 'shows post help without any subcommand' do
      out,err,_ = blogctl("post")
      expect(err).to eq ""
      expect(out).to include("Manage posts")
      expect(out).to include 'Usage'
      expect(out).to include 'Available Commands'
      expect(out).to include 'Use "blogctl post [command] --help" for more information about a command.'

    end
  end

  context 'Basic test cases' do

    let (:blog_path) { "/tmp/blog#{rand 0..10_0000}" }
    
    after(:each) do
      FileUtils.rm_rf(blog_path)
    end

    it 'initializes a new blog' do
      _, err, _ = blogctl("init")
      expect(err).to eq ""
      expect(File.exist?("blog")).to be_truthy
    ensure
      FileUtils.rm_rf("blog")
    end

    it 'initializes a new blog with custom path' do
      _, err, _ = blogctl("init -p #{blog_path}")
      expect(err).to eq ""
      expect(File.exist?(blog_path)).to be_truthy
    end

    it 'fails to initializes a blog twice' do
      out,_,_ = blogctl("init -p #{blog_path}")
      expect(out).not_to match(/Blog environment already exists/)
      expect(File.exist?(blog_path)).to be_truthy
      
      out, _, _ = blogctl("init -p #{blog_path}")
      expect(out).to match(/Blog environment already exists/)
    end
  end

  context 'Post test cases' do

    let (:blog_path) { "/tmp/blog#{rand 0..10_0000}" }
    
    before(:each) do
      blogctl("init -p #{blog_path}")
    end

    after(:each) do
      FileUtils.rm_rf(blog_path)
    end

    it 'fails to add a post without title' do
      out, _, _ = blogctl("post add -p #{blog_path}")
      expect(out).to match(/Title missing/)
    end

    it 'adds a post' do
      out, _, _ = blogctl("post add --title 'test' -p #{blog_path}")
      expect(out).to include "#{blog_path}/test/content.md"
    end
  end
end
